package core

// Define all your config in config.yml and add the model here

// Configurations exported
type Configurations struct {
	Database     DatabaseConfigurations
	EXAMPLE_PATH string
	OTP          OTP
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
	OTP_LENGTH int
}
