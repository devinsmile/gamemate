package gameOwnerController

import (
	"errors"
	"fmt"
	"net/http"

	"sanino/gamemate/controllers/developer"
	"sanino/gamemate/models/game_owner/requests"
	"sanino/gamemate/models/game_owner/responses"
	"sanino/gamemate/models/shared/responses/errors"

	"github.com/labstack/echo"
)

//HandleAddGame handles a request to add a developer API Token.
func HandleAddGame(context echo.Context) error {
	request := gameOwnerRequests.AddGame{}
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}
	if val, err := developerController.IsValidAPI_Token(request.API_Token); !val || err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(errors.New("Rejected by the system, requestor not valid"))
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}
	ownerID, err := getOwnerIDFromSessionToken(request.SessionToken)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(fmt.Errorf("%s token rejected by the system, Invalid Session", request.SessionToken))
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}
	gameID, err := addGameInArchives(ownerID, request.GameName, request.GameDescription, request.MatchMaxPlayers)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(fmt.Errorf("Cannot create new API Token, error => %v", err))
		errorResp.FromError(errors.New("Cannot create API Token"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}
	response := gameOwnerResponses.AddGame{}

	response.FromGameID(gameID)
	return context.JSON(http.StatusCreated, response)
}

//HandleRemoveGame handles a request to remove a developer API Token.
func HandleRemoveGame(context echo.Context) error {
	request := gameOwnerRequests.RemoveGame{}
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}

	IsValid, err := developerController.IsValidAPI_Token(request.API_Token)
	if !IsValid || err != nil {
		context.Logger().Print(fmt.Errorf("API Token %s rejected", request.API_Token))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}
	ownerID, err := getOwnerIDFromSessionToken(request.SessionToken)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(fmt.Errorf("%s token rejected by the system, Invalid Session", request.SessionToken))
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}
	err = removeGameFromCache(request.GameID)
	if err != nil {
		context.Logger().Print(fmt.Errorf("Game with ID:%d not removed. Error => %v", request.GameID, err))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Cannot remove API Token"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}

	err = removeGameFromArchives(ownerID, request.GameID)
	if err != nil {
		context.Logger().Print(fmt.Errorf("Game with ID:%d not removed. Error => %v", request.GameID, err))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Cannot remove API Token"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}

	response := gameOwnerResponses.RemoveGame{}
	response.FromGameID(1)
	return context.JSON(http.StatusOK, response)
}

//HandleRegistration handles a request to register a developer.
func HandleRegistration(context echo.Context) error {
	request := gameOwnerRequests.GameOwnerRegistration{}
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}

	IsValid, err := developerController.IsValidAPI_Token(request.API_Token)
	if !IsValid || err != nil {
		context.Logger().Print(fmt.Errorf("API Token %s rejected", request.API_Token))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}

	ownerID, err := registerOwner(request)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(err, http.StatusInternalServerError)
		return context.JSON(http.StatusBadRequest, errorResp)
	}

	token, err := updateCacheWithSessionOwnerToken(ownerID)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("User registered, but I did not login automatically, try to login later"), http.StatusBadRequest)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}

	responseFromServer := gameOwnerResponses.GameOwnerAuth{}
	responseFromServer.FromToken(token)
	return context.JSON(http.StatusCreated, responseFromServer)
}

//HandleLogin handles login requests for developers.
func HandleLogin(context echo.Context) error {
	request := gameOwnerRequests.GameOwnerAuth{}
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}

	IsValid, err := developerController.IsValidAPI_Token(request.API_Token)
	if !IsValid || err != nil {
		context.Logger().Print(fmt.Errorf("API Token %s rejected", request.API_Token))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}

	isLoggable, ownerID, err := checkLogin(request)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(errors.New("Login failed"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}
	if !isLoggable {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(errors.New("User - Password combination wrong, retry"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}
	token, err := updateCacheWithSessionOwnerToken(ownerID)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(errors.New("Temporary error, retry in a few seconds"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}
	response := gameOwnerResponses.GameOwnerAuth{}
	response.FromToken(token)
	return context.JSON(http.StatusCreated, response)
}