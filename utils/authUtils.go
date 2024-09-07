package utils

import (
	"lamhat/core"
	customErrors "lamhat/errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
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

// Compare OTP
func ValidateOTP(otp_gen_time time.Time, otp_db string, otp_request string) error {

	// Check OTP expired or not
	current_time := time.Now()
	diff_in_mnts := current_time.Sub(otp_gen_time).Minutes()
	if diff_in_mnts > float64(core.Config.OTP.OTP_VALIDITY_IN_MINS) {
		return customErrors.OtpExpired
	}

	// check OTP
	if otp_db == otp_request {
		return nil
	}
	// Wrong OTP
	return customErrors.WrongOtp
}

func GenerateAuthToken(user_id int) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Subject":   user_id,                                                                          // Subject (user identifier)
		"Issuer":    "lamhat-auth",                                                                    // Issuer
		"ExpiresAt": time.Now().Add(time.Duration(core.Config.JWT.JWT_EXP_IN_HRS) * time.Hour).Unix(), // Expiration time
		"IssuedAt":  time.Now().Unix(),                                                                // Issued at
	})

	// Sign the token using the secret key
	tokenString, err := claims.SignedString([]byte(core.Config.JWT.SECRET_KEY))
	if err != nil {
		// Handle error if any
		core.Sugar.Errorf("Failed to sign token: %v", err)
	}
	return tokenString, nil
}
