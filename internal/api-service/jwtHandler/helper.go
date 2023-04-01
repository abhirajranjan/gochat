package jwtHandler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/abhirajranjan/gochat/internal/api-service/serviceErrors"
)

// performs domain logic checks on login request data
//
// returns true if valid login request is created and false otherwise
func validateLoginRequest(request ILoginRequest) serviceErrors.IErr {
	ErrorArray := serviceErrors.ErrorArray{}

	if request.GetUsername() == "" || !IsAlphanum(request.GetUsername()) {
		err := serviceErrors.NewValidationErr("username", "username should be non empty and alpha numeric only")
		ErrorArray = append(ErrorArray, err)
	}

	if request.GetPassword() == "" || !IsAlphanumWithSpecialChar(request.GetPassword()) {
		err := serviceErrors.NewValidationErr("password", fmt.Sprintf("password should be non empty and alpha numeric with %s", SpecialCharacters))
		ErrorArray = append(ErrorArray, err)
	}

	if len(ErrorArray) == 0 {
		return nil
	}
	return &ErrorArray
}

// generate login request from gin.Context
//
// arguments :
// - *gin.Context, has a models.LoginRequest
//
// returns :
// - loginrequest, if validation matches else write StatusBadRequest and returns nil
//
// check only binding errors
func generateLoginRequest(c *gin.Context) (*LoginRequest, error) {
	var request LoginRequest
	err := c.ShouldBind(&request)

	if err != nil {
		return nil, serviceErrors.NewBindingErr(err.Error())
	}

	return &request, nil
}

func checkIfUserHasPermission(permission interface{ Has(string) bool }, reqperm []string) bool {
	for _, perm := range reqperm {
		if !permission.Has(perm) {
			return false
		}
	}
	return true
}

func validateUserData(userdata IUserData) error {
	switch userdata.GetErrCode() {

	case http.StatusOK:
		return nil

	case http.StatusInternalServerError:
		return serviceErrors.ErrInternalServer

	case http.StatusUnauthorized:
		return serviceErrors.NewStandardErr("Status", "invalid credentials")

	default:
		return serviceErrors.ErrInternalServer
	}
}
