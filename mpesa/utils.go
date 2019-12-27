package mpesa
import (
	"fmt"
	"github.com/nyaruka/phonenumbers"
)

// FormatPhoneNumber returns phone number is specific format (E164/National)
func FormatPhoneNumber(phonenumber string, format string) (phone string, err error) {
	num, err := phonenumbers.Parse(phonenumber, "KE")
	if err != nil {
		return
	}
	if ok := phonenumbers.IsValidNumber(num); !ok {
		err = fmt.Errorf("Invalid Phone Number")
		return
	}
	if format == "E164" {
		phone = phonenumbers.Format(num, phonenumbers.E164)
		return
	} else if format == "national" {
		phone = phonenumbers.Format(num, phonenumbers.NATIONAL)
		return
	}
	err = fmt.Errorf("Phone number Format are: NATIONAL, E164")
	return
}
