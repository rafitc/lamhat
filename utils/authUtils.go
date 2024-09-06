package utils

import (
	"lamhat/core"
	customErrors "lamhat/errors"

	"github.com/nanorand/nanorand"
)

// Generate TOTP
func GenerateOtp() (string, error) {
	code, err := nanorand.Gen(core.Config.OTP.OTP_LENGTH)
	if err != nil {
		return "", customErrors.OtpGenError
	}

	return string(code), nil
}
