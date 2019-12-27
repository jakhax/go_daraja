# Go Daraja Api Client
## Description
- Yet another Golang daraja api client library **(WIP)**.

## Contributions
- Highly welcomed, documenting, bug fixes and new features, got some time on you, write the b2b api client!.

## Work Done
- [x] Lipa Na Mpesa Api / Express.
- [ ] C2B Api
- [ ] B2B Api
- [ ] B2C Api
- [ ] Transaction Api
- [ ] Reversal Api

## Installation
```bash
go get github.com/jakhax/go_daraja
# go's google libphonenumber for phone number, i use it for validation
go get github.com/nyaruka/phonenumbers
```

## Usage
- Will document all API's soon, for now refer to the examples below.
- This is not a detailed doc of how daraja api works, for that i highly recommend this article [https://peternjeru.co.ke/safdaraja/ui/](https://peternjeru.co.ke/safdaraja/ui/).

### Lipa Na Mpesa Example
```go
package main;

import (
	"log"
	"fmt"
	"github.com/jakhax/go_daraja"
)

func main(){
	config := &mpesa.Config{
		MpesaEnvironment:"environment e.g sanbox/production",
		MpesaConsumerKey:"Your Consumer Key",
		MpesaConsumerSecret:"Your Consumer Secret",
		MpesaShortCode:"Your Express Short code",
		MpesaExpressPassKey:"Your Express Pass Key",
	}
	expressService,err :=  mpesa.NewExpressService(config);
	if(err != nil){
		log.Fatal(err)
	}
	stkPushRequest, err := expressService.STKPush("PHONE_NUMBER",1,"account",
	"desc","CALLBACK_URL");
	if(err != nil){
		log.Fatal(err)
	}
	fmt.Println(stkPushRequest)
}
```


## References
- [https://peternjeru.co.ke/safdaraja/ui/](https://peternjeru.co.ke/safdaraja/ui/).
