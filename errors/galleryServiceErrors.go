package customErrors

import "errors"

// ErrUserNotFound is a custom error indicating that no user was found.
var ErrStatusNotFound = errors.New("Can't find such status")
var ErrNoGalleryDetailsFound = errors.New("Can't find such gallery")
