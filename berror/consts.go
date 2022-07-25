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
	ErrInternal                   = 1
	ErrBadRequest                 = 2
	ErrForbidden                  = 3
	ErrRPC                        = 4
	ErrCASFailed                  = 5
	ErrNotImplemented             = 6
	ErrKnownBug                   = 7
	ErrUnauthorized               = 8
	ErrUserGeneral                = 101
	ErrUserExisting               = 102
	ErrNotFound                   = 103
	ErrCodeMismatch               = 104
	ErrPasswordMismatch           = 105
	ErrInvalidToken               = 106
	ErrNotExists                  = 107
	ErrUserNameOccupied           = 108
	ErrQueryUserNft               = 109
	ErrUserNftBalanceNotEnough    = 110
	ErrNftAlreadyMintedToUser     = 111
	ErrInvalidState               = 112
	ErrBadTimeRange               = 113
	ErrListingSold                = 114
	ErrListingClosed              = 115
	ErrLockStore                  = 116
	ErrListingBalanceNotEnough    = 117
	ErrListingLocked              = 118
	ErrUserSelfBuy                = 119
	ErrUserSelfOffer              = 120
	ErrUserSelfBid                = 121
	ErrListingTypeNotCorrect      = 122
	ErrBidUnderprice              = 123
	ErrAlreadyListed              = 124
	ErrPaymentMethodNotRegistered = 125
	ErrPaymentInvalid             = 126
	ErrUserCryptoBalanceNotEnough = 127
	ErrTwoFA                      = 128
	ErrExceedMaxPurchaseCount     = 129
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

func NewFail(code int, msg string) *BError {
	return New(code, msg, CategoryBusinessFail)
}
