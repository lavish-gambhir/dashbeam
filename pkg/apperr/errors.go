package apperr

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Original error   `json:"-"`
	Message  string  `json:"message"`
	Code     ErrCode `json:"code"`              // A specific application error code
	Details  any     `json:"details,omitempty"` // Optional additional error details (e.g., validation errors)
}

func (e *Error) Error() string {
	if e.Original != nil {
		return fmt.Sprintf("code: %s, message: %s, original: %s", e.Code, e.Message, e.Original.Error())
	}
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Original
}

func New(code ErrCode, message string) *Error {
	return &Error{
		Message: message,
		Code:    code,
	}
}

func Newf(code ErrCode, format string, args ...any) *Error {
	return &Error{
		Message: fmt.Sprintf(format, args...),
		Code:    code,
	}
}

func Wrap(err error, code ErrCode, message string) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		Original: err,
		Message:  message,
		Code:     code,
	}
}

func Wrapf(err error, code ErrCode, format string, args ...any) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		Original: err,
		Message:  fmt.Sprintf(format, args...),
		Code:     code,
	}
}

func (e *Error) WithDetails(details any) *Error {
	if e != nil {
		e.Details = details
	}
	return e
}

func Is(err error, code ErrCode) bool {
	if err == nil {
		return false
	}
	if ae, ok := err.(*Error); ok {
		return ae.Code == code
	}
	return false
}

func GetCode(err error) ErrCode {
	if err == nil {
		return NoErr
	}
	if ae, ok := err.(*Error); ok {
		return ae.Code
	}
	return Unknown
}

func GetMessage(err error) string {
	if err == nil {
		return ""
	}
	if ae, ok := err.(*Error); ok {
		return ae.Message
	}
	return err.Error()
}

func (e *Error) JSON() []byte {
	b, _ := json.Marshal(e)
	return b
}
