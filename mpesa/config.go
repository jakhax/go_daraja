package mpesa

import (
	"fmt"
)

//Config basic mpesa configurations
type Config struct{
	ConsumerKey string
	ConsumerSecret string
	Environment string 
}

//OK validates config
func (c *Config) OK()(err error){
	if c.ConsumerKey == ""{
		err = fmt.Errorf("ConsumerKey not set")
		return
	}
	if c.ConsumerSecret == ""{
		err = fmt.Errorf("ConsumerSecret not set")
		return
	}
	switch c.Environment {
		case SandBox,Production:
			break
		default:
			err = fmt.Errorf("Invalid Environment options are: sanbox,production")
			return
	}
	return
}
