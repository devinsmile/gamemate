package gameOwnerRequests

import (
	"errors"

	"github.com/labstack/echo"
)

//GameOwnerRegistration represents a request to register a Game Owner into
//the system.
type GameOwnerRegistration struct {
	Type      string `json:"Type" xml:"Type" form:"Type"`
	API_Token string `json:"API_Token" xml:"API_Token" form:"API_Token"`
	Email     string `json:"Email" xml:"Email" form:"Email"`
	Password  string `json:"Password" xml:"Username" form:"Password"`
}

// FromForm Converts from a submitted form (or request) to his struct.
//
// Does not check for the validity of the items inside the struct (e.g. tokens).
func (receiver *GameOwnerRegistration) FromForm(c echo.Context) error {
	var err error
	receiver.Type = c.FormValue("Type")
	receiver.API_Token = c.FormValue("API_Token")
	receiver.Email = c.FormValue("Email")
	receiver.Password = c.FormValue("Password")

	if receiver.Type != "GameOwnerRegistration" || receiver.Email == "" || receiver.Password == "" || receiver.API_Token == "" {
		err = errors.New("Invalid Form Submitted")
	}
	return err
}
