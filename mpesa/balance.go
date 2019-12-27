package mpesa

import (
	"regexp"
	"fmt"
	
)

//BalanceQueryAPI service interface
type BalanceQueryAPI interface{
	BalanceQuery(balanceQuery *BalanceQuery)(apiRes *APIRes,err error)
}

//BalanceQuery model
type BalanceQuery struct{
	InitiatorUserName string
	InitiatorPassword string
	ShortCode string
	IdentifierType string
	TimeOutCallBackURL string
	ResultCallBackURL string
	Remarks string
}

//OK Validates
func (m *BalanceQuery) OK()(err error){
	//shortcode
	digitMatch := regexp.MustCompile(`^[0-9]+$`)
	if !digitMatch.MatchString(m.ShortCode){
		err = fmt.Errorf("ShortCode must be a valid numeric string")
		return
	}
	//IdentiferType
	switch m.IdentifierType{
	case MSISDNIdentiferType,TillNumberIdentifierType,OrganizationIdentifierType:
		break
	case "":
		m.IdentifierType = OrganizationIdentifierType
		break
	default:
		err = fmt.Errorf("Invalid identifier type")
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

//BalanceQueryPayload api payload
type BalanceQueryPayload struct{
	//Initiator	This is the credential/username used to authenticate the transaction request.
	Initiator string `json:"Initiator"`
	//SecurityCredential	Base64 encoded string of the M-Pesa short code and password, 
	//which is encrypted using M-Pesa public key and validates 
	//the transaction on M-Pesa Core system.
	SecurityCredential string `json:"SecurityCredential"`
	//CommandID	A unique command passed to the M-Pesa system.
	CommandID string `json:"CommandID"`
	//PartyA The shortcode of the organisation receiving the transaction.
	PartyA string `json:"PartyA"`
	//IdentifierType Type of the organisation receiving the transaction.
	IdentifierType string `json:"IdentifierType"`
	//Remarks Comments that are sent along with the transaction.
	Remarks string `json:"Remarks"`
	// QueueTimeOutURL	The timeout end-point that receives a timeout message.
	QueueTimeOutURL string `json:"QueueTimeOutURL"`
	// ResultURL	The end-point that receives a successful transaction.
	ResultURL string `json:"ResultURL"`
}

//BalanceQueryRes response 
type BalanceQueryRes struct{

}

//BalanceQuery retur
func (s *Mpesa) BalanceQuery(balanceQuery *BalanceQuery)(apiRes *APIRes,err error){
	err = balanceQuery.OK()
	if err != nil{
		return
	}
	//encrypt password
	securityCredential,err := EncryptPassword(balanceQuery.InitiatorPassword,s.Config.Environment)
	if err != nil{
		return
	}
	payload := &BalanceQueryPayload{
		Initiator:balanceQuery.InitiatorUserName,
		SecurityCredential:securityCredential,
		PartyA:balanceQuery.ShortCode,
		CommandID:AccountBalance,
		IdentifierType:balanceQuery.IdentifierType,
		Remarks:balanceQuery.Remarks,
		QueueTimeOutURL:balanceQuery.TimeOutCallBackURL,
		ResultURL:balanceQuery.ResultCallBackURL,
	} 
	endpoint := "/mpesa/accountbalance/v1/query"
	apiRes, err = s.APIRes(endpoint,payload)
	return
}