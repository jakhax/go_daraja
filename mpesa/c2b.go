package mpesa

import (
	"fmt"
	"regexp"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"bytes"
)

//C2BAPI service  interface
type C2BAPI interface{

}

//RegisterURLs register validation & confirmation url
type RegisterURLs struct{
	ValidationURL string `json:"ValidationURL"`
	ConfirmationURL string `json:"ConfirmationURL"`
	ResponseType string `json:"ResponseType"`
	ShortCode string `json:"ShortCode"`
}

//response types

//CompletedResponseType response type
const CompletedResponseType string = "Completed"

//CancelResponseType response type
const CancelResponseType string = "Cancelled"

//OK validates
func (m *RegisterURLs) OK()(err error){
	//response types
	switch m.ResponseType{
	case CancelResponseType,CompletedResponseType:
		break
	default:
		err = fmt.Errorf("Invalid response type")
		return
	}
	//shortcode
	digitMatch := regexp.MustCompile(`^[0-9]+$`)
	if !digitMatch.MatchString(m.ShortCode){
		err = fmt.Errorf("ShortCode must be a valid numeric string")
		return
	}
	if m.ConfirmationURL == "" && m.ValidationURL==""{
		err = fmt.Errorf("Must provide atlead validation/confirmation url or both")
	}
	return
}



// RegisterURLRes response
// the mis spelled OriginatorCoversationID hence we cannot use apiRes
type RegisterURLRes struct{
	APIRes
	OriginatorCoversationID string `json:"OriginatorCoversationID"`
}


//RegisterURLs register validation and confirmation urls
func (s *Mpesa) RegisterURLs(r *RegisterURLs)(res *RegisterURLRes, err error){
	err = r.OK()
	if err != nil{
		return
	}
	payload,err := json.Marshal(r)
	if err != nil{
		return
	}
	url, err :=  s.GetBaseURL()
	if err != nil{
		return
	}
	url += "/mpesa/c2b/v1/registerurl"
	req, err :=  http.NewRequest(http.MethodPost,url,bytes.NewReader(payload))
	if err != nil{
		return
	}
	req.Header.Add("Content-Type","application/json")
	resp,err := s.MakeRequest(req)
	if err != nil{
		return
	}
	rb,err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil{
		return
	}
	fmt.Println(string(rb))
	if resp.StatusCode != 200 {
		err = s.GetAPIError(resp.Status,resp.StatusCode,rb)
		return
	}
	res = &RegisterURLRes{}
	err = json.Unmarshal(rb,res)
	
	return
}
