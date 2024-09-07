package core

// Define all your config in config.yml and add the model here

// Configurations exported
type Configurations struct {
	Database DatabaseConfigurations
	OTP      OTP
	EMAIL    EMAIL
	JWT      JWT_AUTH
}

// DatabaseConfigurations exported
type DatabaseConfigurations struct {
	HOST        string
	PORT        int
	DATABASE    string
	USER_NAME   string
	PASSWORD    string
	SSL_MODE    string
	CLIENT_CERT string
	CLIENT_KEY  string
	CA_CERT     string
}

type OTP struct {
	OTP_LENGTH           int
	OTP_VALIDITY_IN_MINS int
}

type EMAIL struct {
	MAILGUN_DOMAIN string
	MAILGUN_API    string
	MAINGUN_SENDER string
}

type JWT_AUTH struct {
	JWT_EXP_IN_HRS int
	SECRET_KEY     string
}
