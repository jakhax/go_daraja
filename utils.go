package mpesa

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/nyaruka/phonenumbers"
	"io/ioutil"
	"net/http"
	"time"
)

// FormatPhoneNumber returns phone number is specific format (E164/National)
func FormatPhoneNumber(phonenumber string, format string) (phone string, err error) {
	num, err := phonenumbers.Parse(phonenumber, "KE")
	if err != nil {
		return
	}
	if ok := phonenumbers.IsValidNumber(num); !ok {
		err = &PhoneNumberValidationError{Message: "Invalid Phone Number"}
		return
	}
	if format == "E164" {
		phone = phonenumbers.Format(num, phonenumbers.E164)
		return
	} else if format == "national" {
		phone = phonenumbers.Format(num, phonenumbers.NATIONAL)
		return
	}
	err = errors.New("Phone number Format are: NATIONAL, E164")
	return
}

// AuthToken is the authorization token provided by the daraja api through basic auth
type AuthToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

// NewAuthToken generates a new token via basic auth using consumer key and consumer secret
func NewAuthToken(mc *Config) (authToken *AuthToken, err error) {
	endpoint := "/oauth/v1/generate?grant_type=client_credentials"
	consumerKey, err := mc.GetConsumerkey()
	if err != nil {
		return
	}
	consumerSecret, err := mc.GetConsumerSecret()
	if err != nil {
		return
	}
	password := consumerKey + ":" + consumerSecret
	b64Password := base64.StdEncoding.EncodeToString([]byte(password))

	baseURL, err := mc.GetBaseURL()
	if err != nil {
		return
	}
	url := baseURL + endpoint
	client := http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", "Basic "+b64Password)
	req.Header.Add("Cache-Control", "no-cache")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	jsonBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	authToken = &AuthToken{}
	err = json.Unmarshal(jsonBody, authToken)
	return

}
