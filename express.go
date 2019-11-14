package mpesa;

import (
	"bytes"
	"time"
	"net/http"
	"io"
	"strconv"
	"io/ioutil"
	"encoding/base64"
	"encoding/json"
)


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

type STKPushRes struct{
	MerchantRequestID string
	CheckoutRequestID string
	ResponseCode string
	ResponseDescription string
	CustomerMessage string
}
type STKCallBackResponse  struct {
	Body struct{
		StkCallback STKCallBack `json:"stkCallback"`
	}
}


type STKCallBack struct{
	MerchantRequestID string
	CheckoutRequestID string
	ResultCode int
	ResultDesc string
	CallbackMetadata *STKCallBackItems;
}

type STKCallBackItems struct{
	Item []struct{
		Name string
		Value interface{}
	}
}


type ParsedSTKCallBackRes struct{
	MerchantRequestID string
	CheckoutRequestID string
	ResultCode int
	ResultDesc string
	Meta struct{
		Amount int 
		MpesaReceiptNumber string
		PhoneNumber string
	}
}

type ExpressServiceInterface interface{
	STKPush(phonenumber string, amount int, 
		accountReference string, transactionDesc string,
		callbackURL string) (STKPushRes, error)
	ParseSTKCallbackRes(stkCallBackRes io.Reader) (ParsedSTKCallBackRes, error)
}

type ExpressService struct{
	config MpesaConfig
}

func (s *ExpressService) STKPush(phonenumber string, amount int, 
	accountReference string, transactionDesc string,
	callbackURL string) (STKPushRes, error){

		t := time.Now();
		layout := "20060102150405";
		timestamp := t.Format(layout);
		expressShortCode := s.config.GetExpressShortCode();
		expressPassKey := s.config.GetExpressPassKey();

		password := base64.StdEncoding.EncodeToString([]byte(expressShortCode+expressPassKey+timestamp));
		
		phoneNumber := FormatPhoneNumber(phonenumber, "E164");
		expressPayload := &ExpressPayload{
			BusinessShortCode:expressShortCode,
			Password:password,
			Timestamp:timestamp,
			TransactionType:"CustomerPayBillOnline",
			Amount:amount,
			PartyA:phoneNumber[1:],
			PartyB:expressShortCode,
			PhoneNumber:phoneNumber[1:],
			CallBackURL:callbackURL,
			AccountReference:accountReference,
			TransactionDesc:transactionDesc,
		};
		jsonPayload, err := json.Marshal(expressPayload);
		FatalError(err);
		
		authToken,err := NewAuthToken(s.config);
		FatalError(err);

		client := http.Client{
			Timeout:time.Second*10,
		};

		apiEndpoint := "/mpesa/stkpush/v1/processrequest"
		url := s.config.GetBaseUrl() + apiEndpoint;
		bytesReader := bytes.NewReader(jsonPayload);
		req, err := http.NewRequest(http.MethodPost,url,bytesReader);
		req.Header.Add("Authorization","Bearer "+authToken.AccessToken);
		req.Header.Add("Content-Type","application/json");
		FatalError(err);
		res, err := client.Do(req);
		FatalError(err);
		expressResponse := &STKPushRes{};
		resBody,err := ioutil.ReadAll(res.Body);
		FatalError(err);
		err = json.Unmarshal(resBody, expressResponse);
		FatalError(err);
		return *expressResponse, nil;
	}

func (s *ExpressService)ParseSTKCallbackRes(stkCallBackRes io.Reader) (ParsedSTKCallBackRes, error){
	data,err := ioutil.ReadAll(stkCallBackRes);
	FatalError(err);

	stkCallBack := STKCallBackResponse{};

	err = json.Unmarshal(data, &stkCallBack);
	FatalError(err);

	parsedStkCallBack := ParsedSTKCallBackRes{
		MerchantRequestID:stkCallBack.Body.StkCallback.MerchantRequestID,
		CheckoutRequestID:stkCallBack.Body.StkCallback.CheckoutRequestID,
		ResultCode:stkCallBack.Body.StkCallback.ResultCode,
		ResultDesc:stkCallBack.Body.StkCallback.ResultDesc,
	}

	if stkCallBack.Body.StkCallback.CallbackMetadata != nil {
		for _, item := range stkCallBack.Body.StkCallback.CallbackMetadata.Item {

			switch(item.Name){
				case "Amount":
					amount,_ := item.Value.(float64)
					parsedStkCallBack.Meta.Amount = int(amount);
					break;
				case "MpesaReceiptNumber":
					receipt,_ := item.Value.(string);
					parsedStkCallBack.Meta.MpesaReceiptNumber = receipt;
					break;
				case "PhoneNumber":
					phoneI,_ := item.Value.(float64);
					phone := strconv.Itoa(int(phoneI));
					parsedStkCallBack.Meta.PhoneNumber = phone;
					break;
				default:
					break;
			}
		}
	}

	return parsedStkCallBack, nil;
}

func NewExpressService(mc MpesaConfig)(ExpressService, error){
	expressService := ExpressService{
		config:mc,
	};
	return expressService, nil;
}