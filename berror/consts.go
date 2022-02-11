package berror

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ErrorCategory int

const (
	CategoryBusinessFail      ErrorCategory = 1
	CategoryBusinessTemporary               = 2
	CategorySystemTemporary                 = 3
)

const (
	ErrInternal   = 1
	ErrBadRequest = 2
	ErrForbidden  = 3
	ErrRPC        = 4
	ErrCASFailed  = 5
)

type StackTracer interface {
	StackTrace() errors.StackTrace
}

type BError struct {
	Code          int
	Msg           string
	ErrorCategory ErrorCategory
	StackTrace    errors.StackTrace
	//InnerError    error
}

func (b *BError) Error() string {
	//if b.InnerError != nil {
	//	return fmt.Sprintf("msg: %s, inner: %s", b.Msg, b.InnerError.Error())
	//} else {
	return fmt.Sprintf("code: %d, cat: %d, msg: %s", b.Code, b.ErrorCategory, b.Msg)
	//}
}

func New(code int, msg string, errorCategory ErrorCategory) *BError {
	b := &BError{
		Code:          code,
		Msg:           msg,
		ErrorCategory: errorCategory,
		StackTrace:    errors.New("").(StackTracer).StackTrace(),
	}
	logrus.
		WithField("code", code).
		WithField("msg", msg).
		Error("error")
	fmt.Printf("%+v\n", b.StackTrace)
	return b
}
