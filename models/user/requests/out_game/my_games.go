package outGameRequests

import (
	"errors"

	"github.com/labstack/echo"
)

//UserGameList represents a request dome (or all) the games available to a single user.
type UserGameList struct {
	Type         string `json:"Type" xml:"Type" form:"Type"`
	API_Token    string `json:"API_Token" xml:"API_Token" form:"API_Token"`
	SessionToken string `json:"SessionToken" xml:"SessionToken" form:"SessionToken"`
}

//FromForm creates a valid Sruct based on form data submitted, or returns error.
//
// Does not check for the validity of the items inside the struct (e.g. tokens)
func (receiver *UserGameList) FromForm(c echo.Context) error {
	err := c.Bind(receiver)
	if err != nil || receiver.Type != "UserGameList" {
		return errors.New("Invalid Form Submitted")
	}
	return nil
}
