package errors

import "fmt"

// Error interface type represents an error condition with an error code, with the nil value representing no error.
type Error interface {
	Code() int
	Error() string
}

type err struct {
	code int
	msg  string
}

// New returns a new Error.
func New(code int, msg string) Error {
	return &err{code: code, msg: msg}
}

// Code returns an error code.
func (e err) Code() int {
	return e.code
}

// Error returns an error message.
func (e err) Error() string {
	return e.msg
}

// ErrVarMissing returns new error: "the required variable $varName is missing".
func ErrVarMissing(varName string) error {
	return fmt.Errorf("the required variable $%s is missing", varName)
}
