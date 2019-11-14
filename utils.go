package mpesa;

import (
	"log"
	"time"
	"errors"
	"net/http"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"github.com/nyaruka/phonenumbers"
)

func FatalError(err error){
	if(err != nil){
		log.Fatal(err);
	}
}

func FormatPhoneNumber(phonenumber string, format string)string{
	num, err := phonenumbers.Parse(phonenumber,"KE");
	FatalError(err);
	if(format=="E164"){
		return phonenumbers.Format(num,phonenumbers.E164);
	}else if(format=="national"){
		return phonenumbers.Format(num,phonenumbers.NATIONAL);
	}
	FatalError(errors.New("Phone number Format are: national, E164"));
	return "";
}

type AuthToken struct{
	AccessToken string `json:"access_token"`
	ExpiresIn string `json:"expires_in"`
}

func NewAuthToken(mc MpesaConfig)(*AuthToken, error){
	endpoint := "/oauth/v1/generate?grant_type=client_credentials";
	password := mc.GetConsumerkey()+ ":" +mc.GetConsumerSecret();
	b64Password :=  base64.StdEncoding.EncodeToString([]byte(password));
	
	url := mc.GetBaseUrl()+endpoint;
	client := http.Client{
		Timeout: time.Second*10,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil);
	FatalError(err);
	req.Header.Add("Authorization","Basic "+b64Password);
	req.Header.Add("Cache-Control", "no-cache");

	res, err :=  client.Do(req);
	FatalError(err);
	jsonBody , err := ioutil.ReadAll(res.Body);
	FatalError(err);
	authToken := &AuthToken{};
	err = json.Unmarshal(jsonBody,authToken);
	FatalError(err);
	return authToken, nil;
}