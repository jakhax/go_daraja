package mpesa

//B2CPayload api payload
type B2CPayload struct {
	//username of the API operator in the M-Pesa Org Portal
	InitiatorName string `json:"InitiatorName"`
	//password of the API operator encrypted using the public key certificate provided.
	SecurityCredential string `json:"SecurityCredential"`
	CommandID          string `json:"CommandID"`
	Amount             string `json:"Amount"`
	PartyA             string `json:"Amount"`
	PartyB             string `json:"PartyB"`
	Remarks            string `json:"Remarks"`
	QueueTimeOutURL    string `json:"QueueTimeOutURL"`
	ResultURL          string `json:"ResultURL"`
	Occassion          string `json:"Occassion"`
}

//B2CRes response payload
type B2CRes struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

//B2CService b2c service
type B2CService struct {
}
