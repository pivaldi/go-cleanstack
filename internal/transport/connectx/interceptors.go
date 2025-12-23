package connectx

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/pivaldi/go-cleanstack/internal/platform/apperr"
	"github.com/pivaldi/go-cleanstack/internal/platform/reqid"
	"go.uber.org/zap"
)

type Interceptors struct {
	Log *zap.Logger
}

func (i Interceptors) All() []connect.Interceptor {
	return []connect.Interceptor{
		requestIDInterceptor{i.Log},
		loggingInterceptor{i.Log},
		errorHeaderInterceptor{},
	}
}

type requestIDInterceptor struct{ log *zap.Logger }

func (in requestIDInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		rid := req.Header().Get("X-Request-Id")
		if rid == "" {
			rid = reqid.New()
		}
		req.Header().Set("X-Request-Id", rid) // ensure downstream sees it

		return next(reqid.With(ctx, rid), req)
	}
}
func (in requestIDInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}
func (in requestIDInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}

type errorHeaderInterceptor struct{}

func (in errorHeaderInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		res, err := next(ctx, req)
		if err == nil {
			return res, nil
		}

		// We can’t mutate response headers on error without a response.
		return res, err
	}
}

func (in errorHeaderInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}
func (in errorHeaderInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}

type loggingInterceptor struct{ log *zap.Logger }

func (in loggingInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		start := time.Now()
		res, err := next(ctx, req)
		dur := time.Since(start)

		rid := reqid.Get(ctx)
		fields := []zap.Field{
			zap.String("request_id", rid),
			zap.String("procedure", req.Spec().Procedure),
			zap.String("protocol", req.Peer().Protocol),
			zap.Duration("duration", dur),
			zap.String("peer", req.Peer().Addr),
		}

		if err == nil {
			in.log.Info("rpc", fields...)
			return res, nil
		}

		ae := apperr.As(err)
		if ae == nil {
			in.log.Error("rpc_error", append(fields,
				zap.String("error", err.Error()),
			)...)

			return res, err
		}

		// structured error for Kibana
		fields = append(fields,
			zap.String("err_code", ae.Code),
			zap.String("err_visibility", string(ae.Visibility)),
			zap.Int("http_status", ae.HTTPStatus),
			zap.String("op", ae.Op),
		)

		if ae.Fields != nil {
			fields = append(fields, zap.Any("err_fields", ae.Fields))
		}
		if ae.Req != nil {
			fields = append(fields, zap.Any("req_decoded", ae.Req))
		}
		if ae.Cause != nil {
			fields = append(fields, zap.String("cause", ae.Cause.Error()))
		}
		if ae.Stack != "" {
			fields = append(fields, zap.String("stack", ae.Stack))
		}

		in.log.Error("rpc_error", fields...)

		return res, err
	}
}

func (in loggingInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}
func (in loggingInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}
