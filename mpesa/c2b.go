package mpesa

import (
	"fmt"
	"regexp"
	"encoding/json"
)

//C2BAPI service  interface
type C2BAPI interface{

}


// C2BRes response
// they  mis-spelled OriginatorCoversationID hence we cannot use apiRes
type C2BRes struct{
	APIRes
	OriginatorCoversationID string `json:"OriginatorCoversationID"`
}


//C2BSimulate c2b simulate payment model
type C2BSimulate struct{
	ShortCode string `json:"ShortCode"`
	//optional if not provided will default to CustomerPayBillOnline
	CommandID string `json:"CommandID"`
	Amount float32 `json:"Amount"`
	Msisdn string `json:"Msisdn"`
	BillRefNumber string `json:"BillRefNumber"`
}

//OK validates
func (m *C2BSimulate) OK()(err error){
	//response types
	switch m.CommandID{
	case CustomerPayBillOnline,CustomerBuyGoodsOnline:
		break
	case "":
		m.CommandID = CustomerPayBillOnline
		break
	default:
		err = fmt.Errorf("Invalid response type")
		return
	}
	if m.Msisdn == ""{
		err = fmt.Errorf("Must provide Msisdn")
		return
	}
	phoneNumber, err := FormatPhoneNumber(m.Msisdn,"E164")
	if err != nil{
		return
	}
	//slice +
	m.Msisdn = phoneNumber[1:]

	if m.BillRefNumber == ""{
		m.BillRefNumber = "account"
	}

	//shortcode
	digitMatch := regexp.MustCompile(`^[0-9]+$`)
	if !digitMatch.MatchString(m.ShortCode){
		err = fmt.Errorf("ShortCode must be a valid numeric string")
		return
	}
	return
}

//C2BSimulate simulate c2b payment
func (s *Mpesa) C2BSimulate(c2bSimulate *C2BSimulate)(c2bRes *C2BRes , err error){
	err = c2bSimulate.OK()
	if err != nil{
		return
	}
	endpoint := "/mpesa/c2b/v1/simulate"
	c2bRes, err = s.C2BRes(endpoint,c2bSimulate)
	return

}


//RegisterURLs register validation & confirmation url
type RegisterURLs struct{
	ValidationURL string `json:"ValidationURL"`
	ConfirmationURL string `json:"ConfirmationURL"`
	ResponseType string `json:"ResponseType"`
	ShortCode string `json:"ShortCode"`
}

//response types

//CompletedResponseType response type
const CompletedResponseType string = "Completed"

//CancelResponseType response type
const CancelResponseType string = "Cancelled"

//OK validates
func (m *RegisterURLs) OK()(err error){
	//response types
	switch m.ResponseType{
	case CancelResponseType,CompletedResponseType:
		break
	default:
		err = fmt.Errorf("Invalid response type")
		return
	}
	//shortcode
	digitMatch := regexp.MustCompile(`^[0-9]+$`)
	if !digitMatch.MatchString(m.ShortCode){
		err = fmt.Errorf("ShortCode must be a valid numeric string")
		return
	}
	if m.ConfirmationURL == "" && m.ValidationURL==""{
		err = fmt.Errorf("Must provide atlead validation/confirmation url or both")
	}
	return
}


//RegisterURLs register validation and confirmation urls
func (s *Mpesa) RegisterURLs(r *RegisterURLs)(c2bRes *C2BRes, err error){
	err = r.OK()
	if err != nil{
		return
	}
	endpoint := "/mpesa/c2b/v1/registerurl"
	c2bRes, err = s.C2BRes(endpoint,r)
	return
}


//C2BRes send api request
// this only exists because they mis-spelled "OriginatorCoversationID" hence we cannot use ApiRes
func (s *Mpesa) C2BRes(endpoint string,payload interface{}) (c2bRes *C2BRes, err error){
	rBody,err :=  s.APIRequest(endpoint,payload)
	if err != nil{
		return
	}
	c2bRes = &C2BRes{}
	err = json.Unmarshal(rBody,c2bRes)
	return

}