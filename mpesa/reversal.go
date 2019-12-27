package mpesa

import (
	"regexp"
	"fmt"
	"math"
)

//ReversalAPI service interface
type ReversalAPI interface{
	Reverse(r *Reversal)(apiRes *APIRes, err error)
}

//Reversal model
type Reversal struct{
	TransactionID string
	InitiatorUserName string
	InitiatorPassword string
	// provider either shortcode or phone number depending on the receiver of transaction
	ShortCode string
	PhoneNumber string
	receiverParty string
	Amount float32
	//optional defaults to msdin for phone number and organization for shortcode
	RecieverIdentifierType string
	TimeOutCallBackURL string
	ResultCallBackURL string
	Remarks string
}

//OK validates
func (m *Reversal) OK()(err error){
	if m.ShortCode != "" && m.PhoneNumber != ""{
		err = fmt.Errorf("provider either shortcode or phone number, not both, depending on receiver")
		return
	}

	//shortcode
	if m.ShortCode != ""{
		digitMatch := regexp.MustCompile(`^[0-9]+$`)
		if !digitMatch.MatchString(m.ShortCode){
			err = fmt.Errorf("ShortCode must be a valid numeric string")
			return
		}
		m.receiverParty = m.ShortCode
	}

	if m.PhoneNumber != ""{
		phoneNumber, errX := FormatPhoneNumber(m.PhoneNumber,"E164")
		if errX != nil{
			err = errX
			return
		}
		//slice +
		m.receiverParty = phoneNumber[1:]
		if m.RecieverIdentifierType == ""{
			m.RecieverIdentifierType = MSISDNIdentiferType
		}
	}
	//amount 
	if m.Amount <= float32(0) {
		err = fmt.Errorf("Must provide amount transacted, amount > 0")
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
	//TransactionID
	if m.TransactionID == ""{
		err = fmt.Errorf("Must provide transaction id")
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
	//IdentiferType
	switch m.RecieverIdentifierType{
		case MSISDNIdentiferType,TillNumberIdentifierType,OrganizationIdentifierType:
			break
		case "":
			m.RecieverIdentifierType = OrganizationIdentifierType
			break
		default:
			err = fmt.Errorf("Invalid identifier type")
			return
	}
	return	
}


//ReversalPayload api payload 
type ReversalPayload struct{
	// Initiator	This is the credential/username used to authenticate the transaction request.
	Initiator string `json:"Initiator"`
	// SecurityCredential	Base64 encoded string of the M-Pesa short code and password, 
	//which is encrypted using M-Pesa public key and validates the transaction on M-Pesa Core system.
	SecurityCredential string `json:"SecurityCredential"`
	// CommandID	Unique command for each transaction type, possible values are: TransactionReversal.
	CommandID string `json:"CommandID"`
	// ReceiverParty Organization/MSISDN sending the transaction.
	ReceiverParty string `json:"ReceiverParty"`
	// RecieverIdentifierType	Type of organization receiving the transaction.
	RecieverIdentifierType string `json:"RecieverIdentifierType"`
	// Remarks	Comments that are sent along with the transaction.
	//Amount The amount transacted in that transaction to be reversed, down to the cent.
	Amount float32
	Remarks string `json:"Remarks"`
	// QueueTimeOutURL	The path that stores information of time out transaction.
	QueueTimeOutURL string `json:"QueueTimeOutURL"`
	// ResultURL	The path that stores information of transaction.
	ResultURL string `json:"ResultURL"`
	// TransactionID	Unique identifier to identify a transaction on M-Pesa
	TransactionID string `json:"TransactionID"`
	// Occasion	Optional.
	Occasion string `json:"Occasion"`
}

//Reverse sends request to reverse a transaction
func (s *Mpesa) Reverse(r *Reversal)(apiRes *APIRes, err error){
	err = r.OK()
	if err != nil{
		return
	}
	//encrypt password
	securityCredential,err := EncryptPassword(r.InitiatorPassword,s.Config.Environment)
	if err != nil{
		return
	}
	amount := float32(math.Round(float64(r.Amount)*100)/100)

	payload := &ReversalPayload{
		Initiator:r.InitiatorUserName,
		SecurityCredential:securityCredential,
		ReceiverParty:r.receiverParty,
		CommandID:TransactionReversal,
		RecieverIdentifierType:r.RecieverIdentifierType,
		Remarks:r.Remarks,
		QueueTimeOutURL:r.TimeOutCallBackURL,
		ResultURL:r.ResultCallBackURL,
		TransactionID:r.TransactionID,
		Amount:amount,
	}
	endpoint := "/mpesa/reversal/v1/request"
	apiRes, err = s.SendAPIRequest(endpoint,payload)
	return

}