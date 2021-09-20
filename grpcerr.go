package grpcerr

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error wraps a regular error with a GRPC status code. It can be used just like
// any other error, but when returned by a GRPC serivce endpoint, it will
// resolve to the underlying code. For any other context (e.g. logging), the
// underlying code is not included in the error message.
//
// This makes it possible for a function to return an error with a status code
// that is retained even though the error is wrapped along the way.
type Error struct {
	Code    codes.Code
	Err     error
	Details []proto.Message
}

// New returns a new grpcerr.Error.
func New(code codes.Code, err error) *Error {
	return &Error{
		Code: code,
		Err:  err,
	}
}

// NewMsg returns a new grpcerr.Error by creating an error from the provided
// message.
func NewMsg(code codes.Code, msg string) *Error {
	return &Error{
		Code: code,
		Err:  errors.New(msg),
	}
}

// Error returns the inner error's string without the GRPC status code.
func (e *Error) Error() string {
	return e.Err.Error()
}

// GRPCStatus returns a new GRPC status from the errors' code and error.
func (e *Error) GRPCStatus() *status.Status {
	st := status.New(e.Code, e.Err.Error())
	if detailed, err := st.WithDetails(e.Details...); err != nil {
		return st
	} else {
		return detailed
	}
}

// Unwrap returns the wrapped error. This enables use of errors.Is().
func (e *Error) Unwrap() error {
	return e.Err
}

// Errorf bubbles up GRPC status codes from wrapped errors, allowing for
// composition of rich errors while keeping the code around for the response.
func Errorf(msg string, args ...interface{}) error {
	err := fmt.Errorf(msg, args...)
	code := codes.Unknown
	var details []proto.Message
	if wrapped, ok := err.(interface {
		Unwrap() error
	}); ok {
		if status, ok := status.FromError(wrapped.Unwrap()); ok {
			code = status.Code()
			for _, detail := range status.Details() {
				switch detail := detail.(type) {
				case error:
					continue
				case proto.Message:
					details = append(details, detail)
				}
			}
		}
	}
	return &Error{code, err, details}
}
