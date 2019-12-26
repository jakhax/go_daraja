package mpesa

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// ExpressPayload is the payload for Daraja LNM endpoint
type ExpressPayload struct {
	BusinessShortCode string
	Password          string
	Timestamp         string
	TransactionType   string
	Amount            int
	PartyA            string
	PartyB            string
	PhoneNumber       string
	CallBackURL       string
	AccountReference  string
	TransactionDesc   string
}

// STKPushRes is the response sent by LNM request
type STKPushRes struct {
	MerchantRequestID   string
	CheckoutRequestID   string
	ResponseCode        string
	ResponseDescription string
	CustomerMessage     string
}

// STKCallBackResponse is the response to the LNM callback url
type STKCallBackResponse struct {
	Body struct {
		StkCallback STKCallBack `json:"stkCallback"`
	}
}

// STKCallBack is STKCallBackResponse body
type STKCallBack struct {
	MerchantRequestID string
	CheckoutRequestID string
	ResultCode        int
	ResultDesc        string
	CallbackMetadata  *STKCallBackItems
}

// STKCallBackItems is the array of STKCallBackResponse metadata
type STKCallBackItems struct {
	Item []struct {
		Name  string
		Value interface{}
	}
}

// ParsedSTKCallBackRes is the parsed form of STKCallBackResponse
type ParsedSTKCallBackRes struct {
	MerchantRequestID string
	CheckoutRequestID string
	ResultCode        int
	ResultDesc        string
	Meta              struct {
		Amount             int
		MpesaReceiptNumber string
		PhoneNumber        string
	}
}

// TransactionStatusReq model
type TransactionStatusReq struct {
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	CheckoutRequestID string `json:"CheckoutRequestID"`
}

//TransactionStatusRes api response
type TransactionStatusRes struct {
	MerchantRequestID   string
	CheckoutRequestID   string
	ResponseCode        string
	ResultCode          string
	ResponseDescription string
	ResultDesc          string
	CustomerMessage     string
}

// ExpressServiceInterface interface
type ExpressServiceInterface interface {
	STKPush(phonenumber string, amount int,
		accountReference string, transactionDesc string,
		callbackURL string) (STKPushRes, error)
	ParseSTKCallbackRes(stkCallBackRes io.Reader) (ParsedSTKCallBackRes, error)
}

// ExpressService express api client
type ExpressService struct {
	config *Config
}

// STKPush sends stk push request to daraja
func (s *ExpressService) STKPush(phonenumber string, amount int,
	accountReference string, transactionDesc string,
	callbackURL string) (stkPushRes *STKPushRes, err error) {

	t := time.Now()
	layout := "20060102150405"
	timestamp := t.Format(layout)
	expressShortCode, err := s.config.GetShortCode()
	if err != nil {
		return
	}
	expressPassKey, err := s.config.GetExpressPassKey()
	if err != nil {
		return
	}

	password := base64.StdEncoding.EncodeToString([]byte(expressShortCode + expressPassKey + timestamp))

	phoneNumber, err := FormatPhoneNumber(phonenumber, "E164")
	if err != nil {
		return
	}
	expressPayload := &ExpressPayload{
		BusinessShortCode: expressShortCode,
		Password:          password,
		Timestamp:         timestamp,
		TransactionType:   "CustomerPayBillOnline",
		Amount:            amount,
		PartyA:            phoneNumber[1:],
		PartyB:            expressShortCode,
		PhoneNumber:       phoneNumber[1:],
		CallBackURL:       callbackURL,
		AccountReference:  accountReference,
		TransactionDesc:   transactionDesc,
	}
	jsonPayload, err := json.Marshal(expressPayload)
	if err != nil {
		return
	}

	authToken, err := NewAuthToken(s.config)
	if err != nil {
		return
	}

	client := http.Client{
		Timeout: time.Second * 10,
	}

	apiEndpoint := "/mpesa/stkpush/v1/processrequest"
	url, err := s.config.GetBaseURL()
	if err != nil {
		return
	}
	url = url + apiEndpoint
	bytesReader := bytes.NewReader(jsonPayload)
	req, err := http.NewRequest(http.MethodPost, url, bytesReader)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", "Bearer "+authToken.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	stkPushRes = &STKPushRes{}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(resBody, stkPushRes)
	return
}

// TransactionStatus checks the status of Express Payment
func (s *ExpressService) TransactionStatus(checkOutRequestID string, timeStamp string) (ts *TransactionStatusRes, err error) {
	// validate timestamp
	layout := "20060102150405"
	_, err = time.Parse(layout, timeStamp)
	if err != nil {
		err = errors.New("Invalid Timestamp should be in format: YYYYMMDDHHMMSS")
		return
	}
	expressShortCode, err := s.config.GetShortCode()
	if err != nil {
		return
	}
	expressPassKey, err := s.config.GetExpressPassKey()
	if err != nil {
		return
	}
	password := base64.StdEncoding.EncodeToString([]byte(expressShortCode + expressPassKey + timeStamp))

	transactionStatusReq := &TransactionStatusReq{
		CheckoutRequestID: checkOutRequestID,
		Password:          password,
		Timestamp:         timeStamp,
		BusinessShortCode: expressShortCode,
	}
	data, err := json.Marshal(transactionStatusReq)
	if err != nil {
		return
	}

	authToken, err := NewAuthToken(s.config)
	if err != nil {
		return
	}

	client := http.Client{
		Timeout: time.Second * 10,
	}

	apiEndpoint := "/mpesa/stkpushquery/v1/query"
	url, err := s.config.GetBaseURL()
	if err != nil {
		return
	}
	url = url + apiEndpoint
	bytesReader := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPost, url, bytesReader)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", "Bearer "+authToken.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	body := res.Body
	defer res.Body.Close()
	rBody, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	fmt.Println(string(rBody))
	ts = &TransactionStatusRes{}
	err = json.Unmarshal(rBody, ts)
	return
}

// ParseSTKCallbackRes parses the response from the stk push callback payload
func (s *ExpressService) ParseSTKCallbackRes(stkCallBackRes io.Reader) (parsedStkRes *ParsedSTKCallBackRes,
	err error) {
	data, err := ioutil.ReadAll(stkCallBackRes)
	if err != nil {
		return
	}

	stkCallBack := STKCallBackResponse{}

	err = json.Unmarshal(data, &stkCallBack)
	if err != nil {
		return
	}
	parsedStkRes = &ParsedSTKCallBackRes{
		MerchantRequestID: stkCallBack.Body.StkCallback.MerchantRequestID,
		CheckoutRequestID: stkCallBack.Body.StkCallback.CheckoutRequestID,
		ResultCode:        stkCallBack.Body.StkCallback.ResultCode,
		ResultDesc:        stkCallBack.Body.StkCallback.ResultDesc,
	}

	if stkCallBack.Body.StkCallback.CallbackMetadata != nil {
		for _, item := range stkCallBack.Body.StkCallback.CallbackMetadata.Item {

			switch item.Name {
			case "Amount":
				amount, _ := item.Value.(float64)
				parsedStkRes.Meta.Amount = int(amount)
				break
			case "MpesaReceiptNumber":
				receipt, _ := item.Value.(string)
				parsedStkRes.Meta.MpesaReceiptNumber = receipt
				break
			case "PhoneNumber":
				phoneI, _ := item.Value.(float64)
				phone := strconv.Itoa(int(phoneI))
				parsedStkRes.Meta.PhoneNumber = phone
				break
			default:
				break
			}
		}
	}
	return
}

// NewExpressService returns ExpressService
func NewExpressService(mc *Config) (expressService *ExpressService, err error) {
	expressService = &ExpressService{
		config: mc,
	}
	return
}
