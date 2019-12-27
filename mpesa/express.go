package mpesa

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"regexp"
)

//ExpressAPI service interface
type ExpressAPI interface{
	STKPush(express *Express)(stkPushRes *STKPushRes,err error)
	ParseSTKCallBackRes(stkCallBackRes io.Reader) (parsedStkRes *ParsedSTKCallBackRes,err error)
	ExpressTransactionStatus(shortCode,password,checkOutRequestID string) (ts *ExpressTransactionStatusRes, err error)
}


//Express model 
type Express struct{
	ShortCode string 
	Password string
	TransactionType string
	PhoneNumber string
	CallBackURL string 
	Amount int
	//AccountRef optional defaults to account
	AccountRef string
	//TransactionDesc optional defaults to ""
	TransactionDesc string
}

//OK validates Express model
func (m *Express) OK()(err error){
	//validate shortcode
	digitCheck := regexp.MustCompile(`^[0-9]+$`)
	if !digitCheck.MatchString(m.ShortCode){
		err = fmt.Errorf("ShortCode must be a valid numeric string")
		return
	}
	//password
	if m.Password == ""{
		err = fmt.Errorf("Must provide lnm password")
		return
	}
	//transaction type
	switch m.TransactionType{
		case "":
			//default to paybill
			m.TransactionType = CustomerPayBillOnline
			break
		case CustomerPayBillOnline:
			break
		default:
			err = fmt.Errorf("Invalid transaction type")
			return
	}
	phoneNumber, err := FormatPhoneNumber(m.PhoneNumber,"E164")
	if err != nil{
		return
	}
	// skip +
	m.PhoneNumber = phoneNumber[1:]
	//callbackURl
	if m.CallBackURL == ""{
		err = fmt.Errorf("Must Provide callBackURL")
		return
	}
	//amount
	if m.Amount < 1 {
		err = fmt.Errorf("Amount must be > 0")
		return
	}
	if m.AccountRef == ""{
		m.AccountRef = "account"
	}
	if m.TransactionDesc == ""{
		m.TransactionDesc = "empty desc"
	}
	return
	
}

// ExpressPayload is the payload for Daraja LNM endpoint
type ExpressPayload struct {
	BusinessShortCode string
	Password          string
	Timestamp         string
	TransactionType   string
	Amount            string
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

//STKPush for express api / Lipa Na Mpesa
func (s *Mpesa) STKPush(express *Express)(stkPushRes *STKPushRes,err error){
	err = express.OK()
	if err != nil{
		return
	}
	//timestamp
	t := time.Now()
	layout := "20060102150405"
	timestamp := t.Format(layout)
	//create bs64 password
	password := base64.StdEncoding.EncodeToString([]byte(express.ShortCode + express.Password + timestamp))
	
	//payload
	expressPayload := &ExpressPayload{
		BusinessShortCode: express.ShortCode,
		Password:          password,
		Timestamp:         timestamp,
		TransactionType:   express.TransactionType,
		Amount:            strconv.Itoa(express.Amount),
		PartyA:            express.PhoneNumber,
		PartyB:            express.ShortCode,
		PhoneNumber:       express.PhoneNumber,
		CallBackURL:       express.CallBackURL,
		AccountReference:  express.AccountRef,
		TransactionDesc:   express.TransactionDesc,
	}
	jsonPayload, err := json.Marshal(expressPayload)
	if err != nil {
		return
	}
	apiEndpoint := "/mpesa/stkpush/v1/processrequest"
	url, err := s.GetBaseURL()
	if err != nil {
		return
	}
	url = url + apiEndpoint
	bytesReader := bytes.NewReader(jsonPayload)
	req, err := http.NewRequest(http.MethodPost, url, bytesReader)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	
	res, err := s.MakeRequest(req)
	
	if err != nil {
		return
	}
	resBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return
	}
	if res.StatusCode != 200 {
		err = s.GetAPIError(res.Status,res.StatusCode,resBody)
		return
	}
	stkPushRes = &STKPushRes{}
	err = json.Unmarshal(resBody, stkPushRes)
	return
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

// ParseSTKCallBackRes parses the response from the stk push callback payload
func (s *Mpesa) ParseSTKCallBackRes(stkCallBackRes io.Reader) (parsedStkRes *ParsedSTKCallBackRes,err error) {
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

// ExpressTransactionStatusReq model
type ExpressTransactionStatusReq struct {
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	CheckoutRequestID string `json:"CheckoutRequestID"`
}

//ExpressTransactionStatusRes api response
type ExpressTransactionStatusRes struct {
	MerchantRequestID   string
	CheckoutRequestID   string
	ResponseCode        string
	ResultCode          string
	ResponseDescription string
	ResultDesc          string
	CustomerMessage     string
}


// ExpressTransactionStatus checks the status of Express Payment
func (s *Mpesa) ExpressTransactionStatus(shortCode, password,checkOutRequestID string) (ts *ExpressTransactionStatusRes, err error) {
	// timestamp
	t := time.Now()
	layout := "20060102150405"
	timestamp := t.Format(layout)

	password = base64.StdEncoding.EncodeToString([]byte(shortCode + password + timestamp))

	transactionStatusReq := &ExpressTransactionStatusReq{
		CheckoutRequestID: checkOutRequestID,
		Password:          password,
		Timestamp:         timestamp,
		BusinessShortCode: shortCode,
	}

	data, err := json.Marshal(transactionStatusReq)
	if err != nil {
		return
	}

	apiEndpoint := "/mpesa/stkpushquery/v1/query"
	url, err := s.GetBaseURL()
	if err != nil {
		return
	}
	url = url + apiEndpoint
	bytesReader := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPost, url, bytesReader)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := s.MakeRequest(req)
	if err != nil {
		return
	}
	rBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return
	}
	if res.StatusCode != 200 {
		err = s.GetAPIError(res.Status,res.StatusCode,rBody)
		return
	}
	ts = &ExpressTransactionStatusRes{}
	err = json.Unmarshal(rBody, ts)
	return
}