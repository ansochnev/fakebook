package backend

import (
	"encoding/json"
	"fmt"
)

const (
	ErrBadRequest   = 1000
	ErrMissingParam = 1001
	ErrEmptyParam   = 1002
	ErrInvalidParam = 1003

	// Registration
	ErrEmailExists      = 1101
	ErrUsernameExists   = 1102
	ErrPasswordMismatch = 1103
	ErrPasswordWeak     = 1104

	ErrInternal = 2000
)

type Error struct {
	code    int
	wrapped error
	details any
}

func NewError(code int) *Error {
	return &Error{
		code:    code,
		wrapped: fmt.Errorf("E%d", code),
	}
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Error() string {
	return e.wrapped.Error()
}

func (e *Error) Details() any {
	return e.details
}

func (e *Error) WithMessage(format string, a ...any) *Error {
	e.wrapped = fmt.Errorf("%w: %s", e.wrapped, fmt.Sprintf(format, a...))
	return e
}

func (e *Error) Wrap(err error) *Error {
	e.wrapped = fmt.Errorf("%v: %w", e.wrapped, err)
	return e
}

func (e *Error) With(details any) *Error {
	e.details = details
	return e
}

func (e *Error) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"code":    e.code,
		"message": e.wrapped.Error(),
	}
	if e.details != nil {
		m["details"] = e.details
	}
	return json.Marshal(m)
}

type ParamDetails struct {
	Name  string `json:"name"`
	Value any    `json:"value,omitempty"`
}

func MissingParamError(paramName string) *Error {
	return NewError(ErrMissingParam).
		With(&ParamDetails{
			Name: paramName,
		}).WithMessage("missing mandatory %s", paramName)
}

// ----------------------------------------------------------------------------
//                                Unexported
// ----------------------------------------------------------------------------

func emptyParamError(paramName string) *Error {
	return NewError(ErrEmptyParam).
		With(&ParamDetails{
			Name: paramName,
		}).WithMessage("%s must not be empty", paramName)
}

func invalidParamError(paramName string, value any) *Error {
	return NewError(ErrInvalidParam).
		With(&ParamDetails{
			Name:  paramName,
			Value: value,
		}).WithMessage("invalid %s \"%v\"", paramName, value)
}

func emailExistsError() *Error {
	return NewError(ErrEmailExists).
		WithMessage("email already exists")
}

func usernameExistsError() *Error {
	return NewError(ErrUsernameExists).
		WithMessage("username already exists")
}

func passwordMismatchError() *Error {
	return NewError(ErrPasswordMismatch).
		WithMessage("passwords do not match")
}

func passwordWeakError() *Error {
	return NewError(ErrPasswordWeak).
		WithMessage("weak password")
}

func internalError(err error) *Error {
	return NewError(ErrInternal).
		WithMessage("internal error").Wrap(err)
}

func internalErrorf(format string, a ...any) *Error {
	return NewError(ErrInternal).
		WithMessage("internal error").WithMessage(format, a...)
}

func beginTransactionError(err error) *Error {
	return internalErrorf("cannot begin transaction: %v", err)
}

func commitTransactionError(err error) *Error {
	return internalErrorf("cannot commit transaction: %v", err)
}
