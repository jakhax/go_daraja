package mpesa;

import (
	"fmt"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type C2BApi interface {

}

type B2BApi interface {

}

type TransactionQueryApi interface {

}

type BalanceQueryApi interface{

}

//transaction types

//CustomerPayBillOnline transaction type
const CustomerPayBillOnline string = "CustomerPayBillOnline"
	//CustomerBuyGoodsOnline transaction type
const CustomerBuyGoodsOnline string = "CustomerBuyGoodsOnline"

//environment types

//SandBox environment
const SandBox string = "sandbox"
//Production environment
const Production string = "production"

//commandIDs 

//SalaryPayment b2c commandID
const SalaryPayment string  = "SalaryPayment"
//BusinessPayment b2c commandID
const BusinessPayment string = "BusinessPayment"
//PromotionPayment b2c commandID
const PromotionPayment string = "PromotionPayment"

//Mpesa service implements express, b2c, cb2, b2b, reverse, balance query & transaction query
type Mpesa struct{
	Config *Config
}

//GetBaseURL returns base api url base on environment
func (s *Mpesa) GetBaseURL() (url string, err error) {
	env := s.Config.Environment
	switch env{
		case SandBox:
			url = "https://sandbox.safaricom.co.ke"
			return
		case Production:
			url = "https://api.safaricom.co.ke"
			return
	}
	err = fmt.Errorf("Invalid environment")
	return
}

//Authtoken model
type AuthToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

//GetAuthToken returns *AuthToken
func (s *Mpesa) GetAuthToken()(authToken *AuthToken, err error){
	consumerKey := s.Config.ConsumerKey
	consumerSecret := s.Config.ConsumerSecret
	password := consumerKey + ":" + consumerSecret
	b64Password := base64.StdEncoding.EncodeToString([]byte(password))
	baseURL, err := s.GetBaseURL()
	if err != nil {
		return
	}
	endpoint := "/oauth/v1/generate?grant_type=client_credentials"
	url := baseURL + endpoint
	client := http.Client{}
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

//MakeRequest makes an authenticated http request to daraja api
func (s *Mpesa) MakeRequest(req *http.Request)(res *http.Response, err error){
	client := http.Client{}
	authToken,err := s.GetAuthToken()
	if err != nil{
		return
	}
	req.Header.Add("Authorization", "Bearer "+authToken.AccessToken)
	
	return client.Do(req)
}

//APIError error
type APIError struct{
	RequestID string `json:"requestId"`
	ErrorCode string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Status string
	StatusCode int
}
//Error returns error message
func (e *APIError) Error() string{
	if e.ErrorMessage != "" {
		return e.ErrorMessage
	}
	if e.Status != ""{
		return e.Status
	}
	return fmt.Sprintf("Api error, status code:%d",e.StatusCode)
}

//GetAPIError tries to decode api response to APIError
func (s *Mpesa) GetAPIError(status string, statusCode int, errorBody []byte)(apiError *APIError){
	apiError = &APIError{}
	_ = json.Unmarshal(errorBody,apiError)
	apiError.StatusCode = statusCode
	apiError.Status = status
	return
}

//NewMpesa returns *Mpesa service
func NewMpesa(config *Config)(s *Mpesa, err error){
	err = config.OK()
	if err != nil{
		return
	}
	s = &Mpesa{
		Config:config,
	} 
	return
}