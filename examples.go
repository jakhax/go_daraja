package main

import (
	"fmt"
	"github.com/jakhax/go_daraja/mpesa"
	"log"
)

//MpesaConfig config
var MpesaConfig = &mpesa.Config{
	Environment:    mpesa.SandBox,
	ConsumerKey:    "CONSUMER KEY",
	ConsumerSecret: "CONSUMER SECRET",
}

func sTKPushExample() (err error) {
	mpesaService, err := mpesa.NewMpesa(MpesaConfig)
	if err != nil {
		return
	}
	express := &mpesa.Express{
		ShortCode:   "174379",
		Password:    "LNM Password",
		Amount:      1,
		PhoneNumber: "0712345678",
		CallBackURL: "https://callback.com",
	}
	stkRes, err := mpesaService.STKPush(express)
	if err != nil {
		return
	}
	fmt.Println(stkRes.ResponseDescription)
	return
}

func lnmTransactionStatusExample() (err error) {
	mpesaService, err := mpesa.NewMpesa(MpesaConfig)
	if err != nil {
		return
	}
	res, err := mpesaService.ExpressTransactionStatus("shortcode", "LNM Password", "checkout request id")
	if err != nil {
		return
	}
	fmt.Println(res.ResponseDescription)
	return
}

func b2CExample() (err error) {
	mpesaService, err := mpesa.NewMpesa(MpesaConfig)
	if err != nil {
		return
	}
	b2c := &mpesa.B2C{
		ShortCode:         "123456",
		InitiatorUserName: "testapi115",
		InitiatorPassword: "Safaricom007@",
		PhoneNumber:       "254708374149",
		Amount:            100,
		ResultCallBackURL: "https://callback.com/results",
	}
	res, err := mpesaService.B2C(b2c)
	if err != nil {
		return
	}
	fmt.Println(res.ResponseDescription)
	return
}

func balanceQueryExample() (err error) {
	mpesaService, err := mpesa.NewMpesa(MpesaConfig)
	if err != nil {
		return
	}
	balanceQuery := &mpesa.BalanceQuery{
		ShortCode:         "123456",
		InitiatorUserName: "testapi115",
		InitiatorPassword: "Safaricom007@",
		ResultCallBackURL: "https://callback.com/results",
	}
	res, err := mpesaService.BalanceQuery(balanceQuery)
	if err != nil {
		return
	}
	fmt.Println(res.ResponseDescription)
	return
}

func transactionStatusExample() (err error) {
	mpesaService, err := mpesa.NewMpesa(MpesaConfig)
	if err != nil {
		return
	}
	transactionStatus := &mpesa.TransactionStatus{
		PhoneNumber:       "254712345678",
		InitiatorUserName: "testapi115",
		InitiatorPassword: "Safaricom007@",
		ResultCallBackURL: "https://callback.com/",
		TransactionID:     "LKXXXX1234",
	}
	res, err := mpesaService.TransactionStatus(transactionStatus)
	if err != nil {
		return
	}
	fmt.Println(res.ResponseDescription)
	return
}

func reversalExample() (err error) {
	mpesaService, err := mpesa.NewMpesa(MpesaConfig)
	if err != nil {
		return
	}
	reversal := &mpesa.Reversal{
		PhoneNumber:       "254712345678",
		InitiatorUserName: "testapi115",
		InitiatorPassword: "Safaricom007@",
		ResultCallBackURL: "https://callback.com/",
		TransactionID:     "LKXXXX1234",
		Amount:            1.00,
	}
	res, err := mpesaService.Reverse(reversal)
	if err != nil {
		return
	}
	fmt.Println(res.ResponseDescription)
	return
}

func c2BRegisterURLExample() (err error) {
	mpesaService, err := mpesa.NewMpesa(MpesaConfig)
	if err != nil {
		return
	}
	registerURL := &mpesa.RegisterURLs{
		ValidationURL:   "https://callback.com/validation",
		ConfirmationURL: "https://callback.com/confirmation",
		ShortCode:       "123456",
		// Cancelled/Completed , you can use built in const like below
		ResponseType: mpesa.CompletedResponseType,
	}
	res, err := mpesaService.RegisterURLs(registerURL)
	if err != nil {
		return
	}
	fmt.Println(res.ResponseDescription)
	return
}

func c2BSimulateExample() (err error) {
	mpesaService, err := mpesa.NewMpesa(MpesaConfig)
	if err != nil {
		return
	}
	c2bSimulate := &mpesa.C2BSimulate{
		Amount:    100,
		ShortCode: "123456",
		//phone number
		Msisdn: "254712345678",
	}
	res, err := mpesaService.C2BSimulate(c2bSimulate)
	if err != nil {
		return
	}
	fmt.Println(res.ResponseDescription)
	return
}

func main() {
	err := sTKPushExample()
	if err != nil {
		log.Fatal(err)
	}
}
