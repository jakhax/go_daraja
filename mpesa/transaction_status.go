package mpesa

import (
	"fmt"
	"regexp"
)

//TransactionStatusAPI service interface
type TransactionStatusAPI interface {
	TransactionStatus(ts *TransactionStatus) (apiRes *APIRes, err error)
}

// TransactionStatus model
type TransactionStatus struct {
	TransactionID     string
	InitiatorUserName string
	InitiatorPassword string
	// provider either shortcode or phone number depending on the receiver of transaction
	ShortCode          string
	PhoneNumber        string
	partyA             string
	IdentifierType     string
	TimeOutCallBackURL string
	ResultCallBackURL  string
	Remarks            string
}

//OK validates
func (m *TransactionStatus) OK() (err error) {
	if m.ShortCode != "" && m.PhoneNumber != "" {
		err = fmt.Errorf("provider either shortcode or phone number, not both, depending on receiver")
		return
	}

	//shortcode
	if m.ShortCode != "" {
		digitMatch := regexp.MustCompile(`^[0-9]+$`)
		if !digitMatch.MatchString(m.ShortCode) {
			err = fmt.Errorf("ShortCode must be a valid numeric string")
			return
		}
		m.partyA = m.ShortCode
	}

	if m.PhoneNumber != "" {
		phoneNumber, errX := FormatPhoneNumber(m.PhoneNumber, "E164")
		if errX != nil {
			err = errX
			return
		}
		//slice +
		m.partyA = phoneNumber[1:]
	}

	//initiator username
	if m.InitiatorUserName == "" {
		err = fmt.Errorf("Must provide initiator username")
		return
	}
	//initiator password
	if m.InitiatorPassword == "" {
		err = fmt.Errorf("Must provide initiator password")
		return
	}
	//TransactionID
	if m.TransactionID == "" {
		err = fmt.Errorf("Must provide transaction id")
		return
	}

	if m.ResultCallBackURL == "" {
		err = fmt.Errorf("Must provide a result callback url")
		return
	}
	if m.TimeOutCallBackURL == "" {
		m.TimeOutCallBackURL = m.ResultCallBackURL
	}
	if m.Remarks == "" {
		m.Remarks = "empty remarks"
	}
	//IdentiferType
	switch m.IdentifierType {
	case MSISDNIdentiferType, TillNumberIdentifierType, OrganizationIdentifierType:
		break
	case "":
		m.IdentifierType = OrganizationIdentifierType
		break
	default:
		err = fmt.Errorf("Invalid identifier type")
		return
	}
	return
}

// TransactionStatusPayload api payload
type TransactionStatusPayload struct {
	// CommandID	Unique command for each transaction type,
	// possible values are: TransactionStatusQuery.
	CommandID string `json:"CommandID"`
	// PartyA Organization/MSISDN receiving the transaction
	PartyA string `json:"PartyA"`
	// IdentifierType	Type of organization receiving the transaction
	IdentifierType string `json:"IdentifierType"`
	// Remarks	Comments that are sent along with the transaction.
	Remarks string `json:"Remarks"`
	// Initiator	The name of Initiator to initiating the request.
	Initiator string `json:"Initiator"`
	// SecurityCredential	Base64 encoded string of the M-Pesa
	//short code and password, which is encrypted using M-Pesa
	// public key and validates the transaction on M-Pesa Core system.
	SecurityCredential string `json:"SecurityCredential"`
	// QueueTimeOutURL	The path that stores information of time out transaction.
	QueueTimeOutURL string `json:"QueueTimeOutURL"`
	// ResultURL	The path that stores information of transaction.
	ResultURL string `json:"ResultURL"`
	// TransactionID	Unique identifier to identify a transaction on M-Pesa
	TransactionID string `json:"TransactionID"`
	// Occasion	Optional.
	Occasion string `json:"Occasion"`
}

//TransactionStatus get a transaction's status
func (s *Mpesa) TransactionStatus(ts *TransactionStatus) (apiRes *APIRes, err error) {
	err = ts.OK()
	if err != nil {
		return
	}
	//encrypt password
	securityCredential, err := EncryptPassword(ts.InitiatorPassword, s.Config.Environment)
	if err != nil {
		return
	}

	payload := &TransactionStatusPayload{
		Initiator:          ts.InitiatorUserName,
		SecurityCredential: securityCredential,
		PartyA:             ts.partyA,
		CommandID:          TransactionStatusQuery,
		IdentifierType:     ts.IdentifierType,
		Remarks:            ts.Remarks,
		QueueTimeOutURL:    ts.TimeOutCallBackURL,
		ResultURL:          ts.ResultCallBackURL,
		TransactionID:      ts.TransactionID,
	}
	endpoint := "/mpesa/transactionstatus/v1/query"
	apiRes, err = s.APIRes(endpoint, payload)
	return
}
