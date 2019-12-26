package mpesa

// PhoneNumberValidationError is returned by the format phone number util
type PhoneNumberValidationError struct {
	Message string
}

func (e PhoneNumberValidationError) Error() string {
	return e.Message
}

// ConfigNotSetError when config is not set
type ConfigNotSetError struct {
	Config string
}

func (e ConfigNotSetError) Error() string {
	return e.Config + " Not Set"
}

// InvalidMpesaEnvironment options are: sandbox/production
type InvalidMpesaEnvironment struct {
}

func (e InvalidMpesaEnvironment) Error() string {
	return "Invalid Environment set, options are: sandbox/production"
}
