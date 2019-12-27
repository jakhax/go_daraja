# Go Daraja Api Client
## Description
- Yet another Golang daraja api client library **(WIP)**.

## Work Done
- [x] Lipa Na Mpesa Api / Express.
- [x] C2B Regsiter URL & Simulate Payment Api
- [ ] B2B Api
- [x] B2C Api
- [x] Transaction Api
- [x] Balance Query APi
- [x] Reversal Api
- [ ] Parsers for the callback responses(stk callback parser done)

## Installation
```bash
go get github.com/jakhax/go_daraja
```

## APIs
- This section contain code examples for different apis in daraja
- Note this is not a daraja documentation, refer to references in every section below for detailed documentations.

### Creating the Mpesa Service

```go
package main;
import(
	"fmt"
	"log"
	"github.com/jakhax/go_daraja"
)

//MpesaConfig config
var MpesaConfig = &mpesa.Config{
	//'sandbox' / 'production' , you can use built in consts like below
	Environment:mpesa.SandBox,
	ConsumerKey:"CONSUMER KEY",
	ConsumerSecret:"CONSUMER SECRET",
}
mpesaService, err := mpesa.NewMpesa(config)
``` 

### Express / LNM API
#### LNM STK Push
- To send an STK push to a customer phone
```go
func sTKPushExample()(err error){
	mpesaService, err := mpesa.NewMpesa(MpesaConfig)
	if err != nil{
		return
	}
	express := &mpesa.Express{
		ShortCode:"174379",
		Password:"LNM Password",
		Amount:1,
		PhoneNumber:"0712345678",
		CallBackURL:"https://callback.com",

	}
	stkRes, err := mpesaService.STKPush(express)
	if err != nil{
		return
	}
	fmt.Println(stkRes.ResponseDescription)
	return
}
func main(){
	err := sTKPushExample()
	if err != nil{
		log.Fatal(err)
	}
}	

```

- `mpesaService.STKPush` returns  `STKPushRes` which is a pointer to struct containing the api response. For more information refer to the source.
##### Resources
- [Official API Documentation](https://developer.safaricom.co.ke/lipa-na-m-pesa-online/apis/post/stkpush/v1/processrequest)
- [API quickstart](https://developer.safaricom.co.ke/docs#lipa-na-m-pesa-online-payment)
- [SafDaraja Blog](https://peternjeru.co.ke/safdaraja/ui/#lnm_tutorial)

#### LNM Transaction Status
- Get the transaction status of an lipa na mpesa stk push
```go
func lnmTransactionStatusExample()(err error){
	mpesaService, err := mpesa.NewMpesa(MpesaConfig)
	if err != nil{
		return
	}
	res, err := mpesaService.ExpressTransactionStatus("shortcode","LNM Password","checkout request id")
	if err != nil{
		return
	}
	fmt.Println(res.ResponseDescription)
	return
}
```

##### Resources 
 - [https://developer.safaricom.co.ke/lipa-na-m-pesa-online/apis/post/stkpush/v1/processrequest](https://developer.safaricom.co.ke/lipa-na-m-pesa-online/apis/post/stkpush/v1/processrequest)
 - [https://developer.safaricom.co.ke/docs#lipa-na-m-pesa-online-query-request](https://developer.safaricom.co.ke/docs#lipa-na-m-pesa-online-query-request)

### C2B API

#### C2B Register URL
```go
func c2BRegisterURLExample()(err error){
	mpesaService, err := mpesa.NewMpesa(MpesaConfig)
	if err != nil{
		return
	}
	registerURL := &mpesa.RegisterURLs{
		ValidationURL:"https://callback.com/validation",
		ConfirmationURL:"https://callback.com/confirmation",
		ShortCode:"123456",
		// Cancelled/Completed , you can use built in const like below
		ResponseType:mpesa.CompletedResponseType,
	}
	res,err := mpesaService.RegisterURLs(registerURL)
	if err != nil{
		return
	}
	fmt.Println(res.ResponseDescription)
	return
}
```
##### Resources
- [https://peternjeru.co.ke/safdaraja/ui/#register_url_tutorial](https://peternjeru.co.ke/safdaraja/ui/#register_url_tutorial)
- [https://developer.safaricom.co.ke/c2b/apis/post/registerurl](https://developer.safaricom.co.ke/c2b/apis/post/registerurl)
- [https://developer.safaricom.co.ke/docs#c2b-api](https://developer.safaricom.co.ke/docs#c2b-api)

#### C2B Simulate Transaction
```go
func c2BSimulateExample()(err error){
	mpesaService, err := mpesa.NewMpesa(MpesaConfig)
	if err != nil{
		return
	}
	c2bSimulate := &mpesa.C2BSimulate{
		Amount:100,
		ShortCode:"123456",
		//phone number
		Msisdn:"254712345678",
	}
	res,err := mpesaService.C2BSimulate(c2bSimulate)
	if err != nil{
		return
	}
	fmt.Println(res.ResponseDescription)
	return
}
```
##### Resources
- [https://developer.safaricom.co.ke/c2b/apis/post/simulate](https://developer.safaricom.co.ke/c2b/apis/post/simulate)
- [https://peternjeru.co.ke/safdaraja/ui/#c2b_tutorial](https://peternjeru.co.ke/safdaraja/ui/#c2b_tutorial)
 - [https://developer.safaricom.co.ke/docs#c2b-api](https://developer.safaricom.co.ke/docs#c2b-api)

### B2C API
#### B2C Transaction
```go
func b2CExample()(err error){
	mpesaService, err := mpesa.NewMpesa(MpesaConfig)
	if err != nil{
		return
	}
	b2c := &mpesa.B2C{
		ShortCode:"123456",
		InitiatorUserName:"testapi115",
		InitiatorPassword:"Safaricom007@",
		PhoneNumber:"254712345678",
		Amount:100,
		ResultCallBackURL:"https://callback.com/results",
	}
	res,err := mpesaService.B2C(b2c)
	if err != nil{
		return
	}
	fmt.Println(res.ResponseDescription)
	return
}
```

##### References
- [https://developer.safaricom.co.ke/b2c/apis/post/paymentrequest](https://developer.safaricom.co.ke/b2c/apis/post/paymentrequest)
- [https://peternjeru.co.ke/safdaraja/ui/#b2c_tutorial](https://peternjeru.co.ke/safdaraja/ui/#b2c_tutorial)
- [https://developer.safaricom.co.ke/docs#b2c-api](https://developer.safaricom.co.ke/docs#b2c-api)

### Balance Query API
#### balance query
```go
func balanceQueryExample()(err error){
	mpesaService, err := mpesa.NewMpesa(MpesaConfig)
	if err != nil{
		return
	}
	balanceQuery := &mpesa.BalanceQuery{
		ShortCode:"123456",
		InitiatorUserName:"testapi115",
		InitiatorPassword:"Safaricom007@",
		ResultCallBackURL:"https://callback.com/results",
	}
	res,err := mpesaService.BalanceQuery(balanceQuery)
	if err != nil{
		return
	}
	fmt.Println(res.ResponseDescription)
	return
}
```

##### References
- [https://developer.safaricom.co.ke/account-balance/apis/post/query](https://developer.safaricom.co.ke/account-balance/apis/post/query)
- [https://developer.safaricom.co.ke/docs#account-balance-api](https://developer.safaricom.co.ke/docs#account-balance-api)

### Transaction Status API
#### transaction status
```go
func transactionStatusExample()(err error){
	mpesaService, err := mpesa.NewMpesa(MpesaConfig)
	if err != nil{
		return
	}
	transactionStatus := &mpesa.TransactionStatus{
		PhoneNumber:"254712345678",
		InitiatorUserName:"testapi115",
		InitiatorPassword:"Safaricom007@",
		ResultCallBackURL:"https://callback.com/",
		TransactionID:"LKXXXX1234",
	}
	res,err := mpesaService.TransactionStatus(transactionStatus)
	if err != nil{
		return
	}
	fmt.Println(res.ResponseDescription)
	return
}
```

##### References
- [https://developer.safaricom.co.ke/transaction-status/apis/post/query](
	https://developer.safaricom.co.ke/transaction-status/apis/post/query
)
- [https://developer.safaricom.co.ke/docs#transaction-status](https://developer.safaricom.co.ke/docs#transaction-status)

### Reversal API
#### Reversal
```go
func reversalExample()(err error){
	mpesaService, err := mpesa.NewMpesa(MpesaConfig)
	if err != nil{
		return
	}
	reversal := &mpesa.Reversal{
		PhoneNumber:"254712345678",
		InitiatorUserName:"testapi115",
		InitiatorPassword:"Safaricom007@",
		ResultCallBackURL:"https://callback.com/",
		TransactionID:"LKXXXX1234",
		Amount:1.00,
	}
	res,err := mpesaService.Reverse(reversal)
	if err != nil{
		return
	}
	fmt.Println(res.ResponseDescription)
	return
}
```
##### References
- [https://developer.safaricom.co.ke/reversal/apis/post/request](
	https://developer.safaricom.co.ke/reversal/apis/post/request
)
- [https://peternjeru.co.ke/safdaraja/ui/#reversal_tutorial](https://peternjeru.co.ke/safdaraja/ui/#reversal_tutorial)
- [https://developer.safaricom.co.ke/docs#reversal](https://developer.safaricom.co.ke/docs#reversal)

## Contributions
- Highly welcomed, documenting, report bugs, fix bugs and new features, write the b2b api client.

## References
- [https://developer.safaricom.co.ke/apis-explorer](https://developer.safaricom.co.ke/apis-explorer)
- [https://peternjeru.co.ke/safdaraja/ui/](https://peternjeru.co.ke/safdaraja/ui/)
- [https://developer.safaricom.co.ke/docs](https://developer.safaricom.co.ke/docs)
