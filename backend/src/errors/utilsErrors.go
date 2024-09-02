package customErrors

import "errors"

// ErrUserNotFound is a custom error indicating that no user was found.
var OtpGenError = errors.New("Error while generating error")
