package customErrors

import "errors"

// ErrUserNotFound is a custom error indicating that no user was found.
var ErrUserNotFound = errors.New("user not found")
var ErrNewUserCreation = errors.New("New User creation failed ")
var LoginOtpUpdationFailed = errors.New("Error while updating OTP in DB")
var AuthTokenGenError = errors.New("Error while generating auth token")
