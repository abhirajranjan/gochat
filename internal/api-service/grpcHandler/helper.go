package grpcHandler

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/abhirajranjan/gochat/internal/api-service/model"
	"github.com/abhirajranjan/gochat/internal/api-service/proto/loginService"
	"github.com/abhirajranjan/gochat/internal/api-service/serviceErrors"
	"github.com/abhirajranjan/gochat/pkg/logger"
)

const (
	payloadName = "payload"
)

// handle any internal error that are returned by the grpc endpoint
//
// run when response.ErrCode == http.StatusInternalServerError
func handleInternalGrpcEndpointError(err error, logger logger.ILogger, wrap string) {
	logger.Error(errors.Wrap(err, wrap))
}

// convert model login request struct to grpc proto defined login request
func modelLoginReqToGrpcLoginReq(modelLoginRequest model.ILoginRequest) *loginService.LoginRequest {
	grpcLoginRequest := &loginService.LoginRequest{
		Username: modelLoginRequest.GetUsername(),
		Password: modelLoginRequest.GetPassword(),
	}
	return grpcLoginRequest
}

// performs domain logic checks on login request data
//
// returns true if valid login request is created and false otherwise
func validateLoginRequest(request model.ILoginRequest) serviceErrors.IErr {
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

// convert recieved grpc proto login response to model login response
func grpcLoginResToModelRes(grpcLoginRes *loginService.LoginResponse) model.ILoginResponse {
	user := grpcLoginRes.GetUser()
	status := grpcLoginRes.GetStatus()
	return NewLoginResponse(user.GetUserID(), user.GetUserRoles(), status.GetErrCode(), errors.New(status.GetErr()))
}

// generate map[string]interface{} object for any struct recursively
//
// input: any type
//
// return: map[string]interface{} for struct else input
func GenerateMap(A interface{}, out map[string]interface{}) error {
	return _generateMap(reflect.ValueOf(A), out)
}

// generate map[string]interface{} object from reflect.Value recursively
func _generateMap(val reflect.Value, out map[string]interface{}) error {
	// check if given val is a interface or not and if yes them get the underlying dynamic object
	if val.Kind() == reflect.Interface && !val.IsNil() {
		elm := val.Elem()
		// if interface object if pointer or pointer to pointer then also get the underlying object of it
		if elm.Kind() == reflect.Pointer && !elm.IsNil() && elm.Elem().Kind() == reflect.Pointer {
			val = elm
		}
	}
	// if val is pointer then get the underlying object
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return serviceErrors.NewStandardErr("handler._generateMap", "value passed does not implement struct", val)
	}

	// looping over struct fields
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		// pass is used to check if we have to neglect that field value or not
		// we will neglect field if it has a struct tag of "-"
		pass := false
		name := typeField.Name

		// if there were no struct tag of payloadName specified,
		// then we will use field name as key
		if typeField.Tag.Get(payloadName) != "" {
			name = typeField.Tag.Get(payloadName)
			// if we get struct tag as "-" neglect that field
			if name == "-" {
				pass = true
			}
		}

		if pass {
			continue
		}

		// if the value field is not valid or not exported then return
		if !valueField.IsValid() || !typeField.IsExported() {
			continue
		}

		// extracting the underlying field value if it is an interface
		if valueField.Kind() == reflect.Interface && !valueField.IsNil() {
			elm := valueField.Elem()
			if elm.Kind() == reflect.Pointer && !elm.IsNil() && elm.Elem().Kind() == reflect.Pointer {
				valueField = elm
			}
		}

		// extracting the underlying field value if it a pointer type
		if valueField.Kind() == reflect.Pointer {
			valueField = valueField.Elem()

		}

		// field value is itself a struct
		// generate map of it recursively
		if valueField.Kind() == reflect.Struct {
			recursiveout := map[string]interface{}{}
			err := _generateMap(valueField, recursiveout)
			if err != nil {
				return errors.Wrap(err, "handler._generateMap")
			}
			out[name] = recursiveout
			continue
		}

		// actual field value as an interface
		outval := valueField.Interface()
		out[name] = outval
	}
	return nil
}

func GetVersionFromClaims(claims map[string]interface{}) (int64, error) {
	version, ok := claims["ver"]
	if !ok {
		return 0, serviceErrors.NewStandardErr("handler.GetVersionFromClaims", "jwt missing version", claims)

	}
	ver, ok := version.(int64)
	if !ok {
		return 0, serviceErrors.NewStandardErr("handler.ExtractPayloadData", "jwt has unknown type version", version)
	}
	return ver, nil
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
func generateLoginRequest(c *gin.Context) (model.ILoginRequest, error) {
	var request LoginRequest
	err := c.ShouldBind(&request)

	if err != nil {
		return nil, serviceErrors.NewBindingErr(err.Error())
	}

	return &request, nil
}
