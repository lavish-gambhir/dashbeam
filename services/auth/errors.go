package auth

import "github.com/lavish-gambhir/dashbeam/pkg/apperr"

var (
	ErrInvalidCredentials = apperr.New(apperr.InvalidCredentials, "invalid username or password")
	ErrUserNotFound       = apperr.New(apperr.UserNotFound, "user not found")
	ErrUserInactive       = apperr.New(apperr.InvalidCredentials, "user account is inactive")
	ErrTokenExpired       = apperr.New(apperr.TokenExpired, "token has expired")
)
