package mpesa

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

//B2CService b2c service
type B2CService struct {
}
