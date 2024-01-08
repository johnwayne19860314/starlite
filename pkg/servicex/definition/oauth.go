package def

import "errors"

var (
	ErrAPINeedSessionToken = errors.New("API need session token")
	ErrInvalidSessionToken = errors.New("invalid session token")
	ErrInvalidAuthToken    = errors.New("invalid auth token")
	ErrNotFoundAccessToken = errors.New("not found access token")
)
