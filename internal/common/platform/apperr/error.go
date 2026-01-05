package apperr

import (
	"errors"
	"net/http"
	"runtime/debug"
	"time"
)

type Visibility string

const (
	Public  Visibility = "public"
	Private Visibility = "private"
)

type AppError struct {
	Code       string     // stable error code (for clients + logs)
	Message    string     // safe when Public, generic when Private
	HTTPStatus int        // desired HTTP status mapping
	Visibility Visibility // public/private

	Op     string         // operation label (use-case / handler)
	When   int64          // unix timestamp
	Cause  error          // wrapped error (private)
	Stack  string         // stacktrace for private error
	Req    map[string]any // decoded request summary (sanitized)
	Fields map[string]any // arbitrary structured fields
}

func (e *AppError) Error() string   { return e.Code + ": " + e.Message }
func (e *AppError) Unwrap() error   { return e.Cause }
func (e *AppError) IsPrivate() bool { return e.Visibility == Private }
func (e *AppError) IsPublic() bool  { return e.Visibility == Public }

func NewPublic(code, msg string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    msg,
		HTTPStatus: httpStatus,
		Visibility: Public,
		When:       time.Now().Unix(),
	}
}

func WrapPrivate(code string, httpStatus int, cause error) *AppError {
	return &AppError{
		Code:       code,
		Message:    "internal error",
		HTTPStatus: httpStatus,
		Visibility: Private,
		Cause:      cause,
		Stack:      string(debug.Stack()),
		When:       time.Now().Unix(),
	}
}

func As(err error) *AppError {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae
	}

	return nil
}

// Helpers to decorate errors without losing stack.
func WithOp(err error, op string) error {
	if ae := As(err); ae != nil {
		ae.Op = op
	}

	return err
}

func WithReq(err error, req map[string]any) error {
	if ae := As(err); ae != nil {
		ae.Req = req
	}

	return err
}

func WithField(err error, k string, v any) error {
	if ae := As(err); ae != nil {
		if ae.Fields == nil {
			ae.Fields = map[string]any{}
		}
		ae.Fields[k] = v
	}

	return err
}

func StatusOrDefault(ae *AppError, def int) int {
	if ae == nil || ae.HTTPStatus <= 0 {
		return def
	}

	return ae.HTTPStatus
}

// Convenience constructors
func BadRequest(code, msg string) *AppError   { return NewPublic(code, msg, http.StatusBadRequest) }
func NotFound(code, msg string) *AppError     { return NewPublic(code, msg, http.StatusNotFound) }
func Conflict(code, msg string) *AppError     { return NewPublic(code, msg, http.StatusConflict) }
func Unauthorized(code, msg string) *AppError { return NewPublic(code, msg, http.StatusUnauthorized) }
