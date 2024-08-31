package customErrors

import "errors"

// ErrUserNotFound is a custom error indicating that no user was found.
var ErrUserNotFound = errors.New("user not found")
