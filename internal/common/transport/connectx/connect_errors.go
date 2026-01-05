package connectx

import (
	"errors"
	"net/http"

	"connectrpc.com/connect"
	"github.com/pivaldi/go-cleanstack/internal/common/platform/apperr"
)

func ConnectCodeFromHTTPStatus(st int) connect.Code {
	switch st {
	case http.StatusBadRequest:
		return connect.CodeInvalidArgument
	case http.StatusUnauthorized:
		return connect.CodeUnauthenticated
	case http.StatusForbidden:
		return connect.CodePermissionDenied
	case http.StatusNotFound:
		return connect.CodeNotFound
	case http.StatusConflict:
		return connect.CodeAlreadyExists
	case http.StatusTooManyRequests:
		return connect.CodeResourceExhausted
	case http.StatusNotImplemented:
		return connect.CodeUnimplemented
	case http.StatusServiceUnavailable:
		return connect.CodeUnavailable
	default:
		return connect.CodeInternal
	}
}

// ToConnectError returns a Connect error that is safe for clients.
// We also inject the stable error code into response headers via interceptor.
func ToConnectError(err error) error {
	ae := apperr.As(err)
	if ae == nil {
		return connect.NewError(connect.CodeInternal, errors.New("internal error"))
	}
	cc := ConnectCodeFromHTTPStatus(apperr.StatusOrDefault(ae, http.StatusInternalServerError))

	// Public: keep message. Private: never leak.
	msg := ae.Message
	if ae.Visibility == apperr.Private {
		msg = "internal error"
	}

	return connect.NewError(cc, errors.New(msg))
}
