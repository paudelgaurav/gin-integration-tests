package apierrors

import "errors"

var (
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrFailedToGenerateToken = errors.New("failed to generate tokens")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrUnauthorizedAccess    = errors.New("unauthorized access")
	ErrAuthHeaderMissing     = errors.New("authorization header is required")
	ErrInvalidTokenFormat    = errors.New("invalid token format")
	ErrInvalidOrExpiredToken = errors.New("invalid or expired token")
	ErrInvalidTokenClaims    = errors.New("invalid token claims")
	ErrFailedToUpdateUser    = errors.New("failed to update user")
	ErrFailedToDeleteUser    = errors.New("failed to delete user")
	ErrTooManyRequests       = errors.New("too many requests")
)
