package mpesa

import (
	"fmt"
	"regexp"
	"strconv"
)

//B2CAPI service intercface
type B2CAPI interface {
	B2C(b2c *B2C) (apiRes *APIRes, err error)
}

//B2C model 
type B2C struct{
	InitiatorUserName string
	InitiatorPassword string
	ShortCode string
	PhoneNumber string
	Amount int 
	//optional defaults to BusinessPayment
	CommandID string
	ResultCallBackURL string
	//optional defaults to ResultURL
	TimeOutCallBackURL string
	//optional defaults to ""
	Remarks string
}

//OK validates B2C
func (m *B2C) OK() (err error){
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
	if m.InitiatorPassword == ""{
		err= fmt.Errorf("Must provide initiator password")
		return
	}

	//phonenumber
	if m.PhoneNumber == ""{
		err= fmt.Errorf("Must provide phone number")
		return
	}
	phoneNumber, err := FormatPhoneNumber(m.PhoneNumber,"E164")
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
	if m.TimeOutCallBackURL == ""{
		m.TimeOutCallBackURL = m.ResultCallBackURL
	}
	if m.Remarks == ""{
		m.Remarks = "empty remarks"
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


//B2C sends a b2c request to daraja
func (s *Mpesa) B2C(b2c *B2C)(apiRes *APIRes, err error){
	err = b2c.OK()
	if err != nil{
		return
	}
	//encrypt password
	securityCredential,err := EncryptPassword(b2c.InitiatorPassword,s.Config.Environment)
	if err != nil{
		return
	}
	payload := &B2CPayload{
		InitiatorName:b2c.InitiatorUserName,
		SecurityCredential:securityCredential,
		CommandID:b2c.CommandID,
		Amount:strconv.Itoa(b2c.Amount),
		PartyA:b2c.ShortCode,
		PartyB:b2c.PhoneNumber,
		Remarks:b2c.Remarks,
		QueueTimeOutURL:b2c.TimeOutCallBackURL,
		ResultURL:b2c.ResultCallBackURL,
	}
	endpoint := "/mpesa/b2c/v1/paymentrequest"
	apiRes, err = s.SendAPIRequest(endpoint,payload)
	return
}
