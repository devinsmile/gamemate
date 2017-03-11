package developerResponses

import "strings"

//AddToken represents a POSITIVE response from the server to a developerRequests.AddToken
//
//For NEGATIVE response, please refer to errorResponses.ErrorResponse.
type AddToken struct {
	Type         string `json:"Type" xml:"Type" form:"Type"`
	NewAPI_Token string `json:"SessionToken" xml:"SessionToken" form:"SessionToken"`
}

//FromAPIToken fills the struct's data with proper definition, based on an
//API token.
func (receiver *AddToken) FromAPIToken(API_Token string) {
	receiver.Type = "AddToken"
	receiver.NewAPI_Token = strings.Replace(receiver.NewAPI_Token, "0x", "", 1)
}
