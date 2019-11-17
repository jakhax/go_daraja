package mpesa

// ConfigInterface : mpesa settings interface
type ConfigInterface interface {
	GetConsumerkey() (c string, err error)
	GetConsumerSecret() (c string, err error)
	GetEnvironment() (c string, err error)
	GetShortCode() (c string, err error)
	// GetExpressShortCode() (c string, err error)
	GetExpressPassKey() (c string, err error)
	GetBaseURL() (c string, err error)
}

// Config : mpesa settings
type Config struct {
	MpesaConsumerKey    string
	MpesaConsumerSecret string
	MpesaEnvironment    string
	MpesaShortCode      string
	// MpesaExpressShortCode string
	MpesaExpressPassKey string
}

// GetConsumerkey returns the consumer key
func (mc *Config) GetConsumerkey() (c string, err error) {
	c = mc.MpesaConsumerKey
	if mc.MpesaConsumerKey == "" {
		err = &ConfigNotSetError{Config: "Mpesa Consumer Key"}
	}
	return
}

// GetConsumerSecret returns the Consumer Secret
func (mc *Config) GetConsumerSecret() (c string, err error) {
	c = mc.MpesaConsumerSecret
	if mc.MpesaConsumerSecret == "" {
		err = &ConfigNotSetError{Config: "Mpesa Consumer Secret"}
	}
	return

}

// GetEnvironment returns enviroment, should be sandbox or production
func (mc *Config) GetEnvironment() (c string, err error) {
	c = mc.MpesaEnvironment
	if mc.MpesaEnvironment == "" {
		err = &ConfigNotSetError{Config: "Mpesa Envrironment Secret"}
	}
	return
}

// GetShortCode returns the short code (till no/ pay bill)
func (mc *Config) GetShortCode() (c string, err error) {
	c = mc.MpesaShortCode
	if mc.MpesaShortCode == "" {
		err = &ConfigNotSetError{Config: "Mpesa Short Code"}
	}
	return
}

// func (mc *Config) GetExpressShortCode() (c string, err error){
// 	c = mc.MpesaExpressShortCode;
// 	if mc.MpesaExpressShortCode == ""{
// 		err = &ConfigNotSetError{Config:"Mpesa Express Short Code Secret"};
// 	}
// 	return;
// }

// GetExpressPassKey returns the Express PassKey / LNM Passkey
func (mc *Config) GetExpressPassKey() (c string, err error) {
	c = mc.MpesaExpressPassKey
	if mc.MpesaExpressPassKey == "" {
		err = &ConfigNotSetError{Config: "Mpesa Express Pass Key"}
	}
	return
}

// GetBaseURL returns base  daraja Api url depending on environment
func (mc *Config) GetBaseURL() (c string, err error) {
	mpesaEnvironment, err := mc.GetEnvironment()
	if err != nil {
		return
	}
	if mpesaEnvironment == "sandbox" {
		c = "https://sandbox.safaricom.co.ke"
	} else if mpesaEnvironment == "production" {
		c = "https://api.safaricom.co.ke"
	} else {
		err = &InvalidMpesaEnvironment{}
	}
	return
}
