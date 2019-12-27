package mpesa

import (
	"fmt"
	"regexp"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"net/http"
	"strconv"
)

type B2CApi interface {
	B2C(b2c *B2C) (b2cRes *B2CRes, err error)
}

//B2C model 
type B2C struct{
	InitiatorUserName string
	Password string
	ShortCode string
	PhoneNumber string
	Amount int 
	//optional defaults to BusinessPayment
	CommandID string
	ResultCallBackURL string
	//optional defaults to ResultURL
	TimeoutCallBackURL string
	//optional defaults to ""
	Remarks string
}

func (m *B2C) OK(err error){
	//shortcode
	digitMatch := regexp.MustCompile(`^[0-9]+$`)
	if !digitMatch.MatchString(m.ShortCode){
		err = fmt.Errorf("ShortCode must be a valid numeric string")
		return
	}
	//initiator username
	if m.InitiatorUserName == ""{
		err = fmt.Errorf("Must provide initiator username")
		return
	}
	//initiator password
	if m.Password == ""{
		err= fmt.Errorf("Must provide initiator password")
		return
	}

	//phonenumber
	if m.PhoneNumber == ""{
		err= fmt.Errorf("Must provide phone number")
		return
	}
	phoneNumber, err = FormatPhoneNumber(m.PhoneNumber,"E164")
	if err != nil{
		return
	}
	//slice +
	m.PhoneNumber = phoneNumber[1:]

	//commandId
	switch m.CommandID {
		case SalaryPayment, BusinessPayment, PromotionPayment:
			break
		case "":
			m.CommandID = BusinessPayment
			break
		default:
			err = fmt.Errorf("Invalid CommandID")
			return
	}
	if m.Amount <= 0{
		err = fmt.Errorf("Amount must be > 0")
		return
	}
	if m.ResultCallBackURL == ""{
		err = fmt.Errorf("Must provide a result callback url")
		return
	}
	if m.TimeoutCallBackURL == ""{
		m.TimeoutCallBackURL = m.ResultCallBackURL
	}
	return
}


//B2CPayload api payload
type B2CPayload struct {
	//InitiatorNames is the credential/username used to
	//authenticate the transaction request.
	InitiatorName string `json:"InitiatorName"`
	//SecurityCredential is the Base64 encoded string of the
	//B2C short code and password, which is encrypted using
	//M-Pesa public key and validates the transaction on M-Pesa Core system.
	SecurityCredential string `json:"SecurityCredential"`
	//CommandID	Unique command for each transaction type
	//e.g. SalaryPayment, BusinessPayment, PromotionPayment
	CommandID string `json:"CommandID"`
	//Amount The amount being transacted
	Amount string `json:"Amount"`
	//PartyA is the Organization’s shortcode initiating the transaction.
	PartyA string `json:"PartyA"`
	//PartyB Phone number receiving the transaction
	PartyB string `json:"PartyB"`
	//Remarks Comments that are sent along with the transaction.
	Remarks string `json:"Remarks"`
	//QueueTimeOutURL The timeout end-point that receives a timeout response.
	QueueTimeOutURL string `json:"QueueTimeOutURL"`
	//ResultURL	The end-point that receives the response of the transaction
	ResultURL string `json:"ResultURL"`
	Occassion string `json:"Occassion"`
}

//B2CRes response payload
type B2CRes struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

//B2C sends a b2c request to daraja
func (s *Mpesa) B2C(b2c *B2C)(b2cRes *B2CRes, err error){
	err = b2c.OK()
	if err != nil{
		return
	}
	//encrypt password
	cipherText,err := EncryptPassword(b2c.Password,s.Config.Environment)
	if err != nil{
		return
	}
	securityCredential = base64.StdEncoding.EncodeToString(cipherText)
	payload := &B2CPayload{
		InitiatorName:b2c.InitiatorUserName,
		SecurityCredential:securityCredential,
		CommandID:b2c.CommandID,
		Amount:strconv.Itoa(b2c.Amount),
		PartyA:b2c.ShortCode,
		PartyB:b2c.PhoneNumber,
		Remarks:b2c.Remarks,
		QueueTimeOutURL:b2c.TimeoutCallBackURL,
		ResultURL:b2c.ResultCallBackURL,
	}
	jsonPayload,err:= json.Marshal(payload)
	if err != nil{
		return
	}
	requestPayload := bytes.NewReader(jsonPayload)
	url,err := s.GetBaseURL()
	if err != nil{
		return
	}
	url += "/mpesa/b2c/v1/paymentrequest"
	req,err := http.NewRequest(http.MethodPost,url,requestPayload)
	req.Header.Add("Content-Type","application/json")
	res,err := s.MakeRequest(req)
	if err != nil{
		return
	}

	body := res.Body 
	rBody, err := ioutil.ReadAll(body)
	defer res.Body.Close()
	if err != nil{
		return
	}
	b2cRes = &B2CRes{}
	err = json.Unmarshal(rBody,b2cRes)
	return
}
