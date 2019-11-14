package mpesa;

import (
	"log"
)

type MpesaConfigInterface interface{
	GetConsumerkey() string
	GetConsumerSecret() string
	GetEnvironment() string
	GetShortCode() string
	GetExpressShortCode() string 
	GetExpressPassKey() string
	GetBaseUrl() string
}

func VariableNotSetFatal(s string){
	log.Fatalf("%s not set\n",s);
}


type MpesaConfig struct {
	MpesaConsumerKey string
	MpesaConsumerSecret string
	MpesaEnvironment string
	MpesaShortCode string
	MpesaExpressShortCode string
	MpesaExpressPassKey string
}


func (mc *MpesaConfig) GetConsumerkey()string{
	if(mc.MpesaConsumerKey == ""){
		VariableNotSetFatal("Mpesa Consumer Key");
	}
	return mc.MpesaConsumerKey;
}

func (mc *MpesaConfig) GetConsumerSecret()string{
	if(mc.MpesaConsumerSecret == ""){
		VariableNotSetFatal("Mpesa Consumer Secret");
	}
	return mc.MpesaConsumerSecret;
}

func (mc *MpesaConfig) GetEnvironment()string{
	if(mc.MpesaEnvironment == ""){
		VariableNotSetFatal("MpesaEnvironment");
	}
	return mc.MpesaEnvironment;
}

func (mc *MpesaConfig) GetShortCode()string{
	if(mc.MpesaShortCode == ""){
		VariableNotSetFatal("Mpesa ShortCode");
	}
	return mc.MpesaShortCode;
}

func (mc *MpesaConfig) GetExpressShortCode()string{
	if(mc.MpesaExpressShortCode == ""){
		VariableNotSetFatal("Mpesa Express ShortCode");
	}
	return mc.MpesaExpressShortCode;
}

func (mc *MpesaConfig) GetExpressPassKey()string{
	if(mc.MpesaExpressPassKey == ""){
		VariableNotSetFatal("Mpesa Express Pass Key");
	}
	return mc.MpesaExpressPassKey;
}

func (mc *MpesaConfig) GetBaseUrl()string{
	mpesaEnvironment := mc.GetEnvironment();
	if(mpesaEnvironment=="sandbox"){
		return "https://sandbox.safaricom.co.ke";
	}else if(mpesaEnvironment=="production"){
		return "https://api.safaricom.co.ke"
	}else{
		log.Fatal("Mpesa Environment options are: sandbox or production");
	}
	return "";
}