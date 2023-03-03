package grpcHandler

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/abhirajranjan/gochat/internal/api-service/proto/loginService"
	"github.com/abhirajranjan/gochat/internal/api-service/serviceErrors"
)

const (
	payloadName = "payload"
)

// convert model login request struct to grpc proto defined login request
func modelLoginReqToGrpcLoginReq(modelLoginRequest ILoginRequest) (grpcLoginRequest *loginService.LoginRequest) {
	grpcLoginRequest.Username = modelLoginRequest.GetUsername()
	grpcLoginRequest.Password = modelLoginRequest.GetPassword()

	return grpcLoginRequest
}

// performs domain logic checks on login request data
//
// returns true if valid login request is created and false otherwise
func validateLoginRequest(request ILoginRequest) serviceErrors.IErr {
	ErrorArray := serviceErrors.ErrorArray{}

	if request.GetUsername() == "" && !IsAlphanum(request.GetUsername()) {
		err := serviceErrors.NewValidationErr("username", "username should be non empty and alpha numeric only")
		ErrorArray = append(ErrorArray, err)
	}

	if request.GetPassword() == "" && !IsAlphanumWithSpecialChar(request.GetPassword()) {
		err := serviceErrors.NewValidationErr("password", fmt.Sprintf("password should be non empty and alpha numeric with %s", SpecialCharacters))
		ErrorArray = append(ErrorArray, err)
	}

	if len(ErrorArray) == 0 {
		return nil
	}
	return &ErrorArray
}

// convert recieved grpc proto login response to model login response
func grpcLoginResToModelRes(grpcLoginRes *loginService.LoginResponse) ILoginResponse {
	user := grpcLoginRes.GetUser()
	status := grpcLoginRes.GetStatus()
	return NewLoginResponse(user.GetUserID(), user.GetUserRoles(), status.GetErrCode(), errors.New(status.GetErr()))
}

// generate map[string]interface{} object for any struct recursively
//
// input: any type
//
// return: map[string]interface{} for struct else input
func GenerateMap(A interface{}) interface{} {
	return _generateMap(reflect.ValueOf(A))
}

// generate map[string]interface{} object from reflect.Value recursively
func _generateMap(val reflect.Value) interface{} {
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

	out := map[string]interface{}{}

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

		// actual field value as an interface
		outval := valueField.Interface()

		// field value is itself a struct
		// generate map of it recursively
		if valueField.Kind() == reflect.Struct {
			outval = _generateMap(valueField)
		}

		out[name] = outval

	}
	return out
}

// extract map[string]interface{} to struct
func ExtractMapTo(maps map[string]interface{}, to any) {
	_extractMapTo(maps, reflect.ValueOf(to))
}

// extract map[string]interface{} to reflect.Value of struct
func _extractMapTo(claims map[string]interface{}, to reflect.Value) {
	// if given reflect.Value is interface then get the underlying object
	if to.Kind() == reflect.Interface && to.IsNil() {
		elm := to.Elem()
		// if interface object is pointer or pointer to pointer then get the object
		if elm.Kind() == reflect.Pointer && !elm.IsNil() && elm.Elem().Kind() == reflect.Pointer {
			to = elm
		}
	}

	// if given reflect.Value is a pointer type then get the underlying object
	if to.Kind() == reflect.Pointer {
		to = to.Elem()
	}

	// the final object type should be struct to decode data into it
	if to.Kind() != reflect.Struct {
		return
	}

	// *looping over all struct field to check if the map provided has data related to it or not
	// *cannot loop over map directly as it will not allow us to match the struct tags as well
	for i := 0; i < to.NumField(); i++ {
		// get the field reflection
		field := to.Field(i)

		// claimval is the the map value for the key if key matches field
		var claimval interface{}

		// check if the given field has a tag and if yes, is it in our map
		tag := to.Type().Field(i).Tag.Get(payloadName)

		var ok bool

		if tag != "" {
			claimval, ok = claims[tag]
		}
		if !ok {
			claimval, ok = claims[to.Type().Field(i).Name]
		}
		if !ok {
			continue
		}

		// check if the field is valid
		if !field.IsValid() {
			continue
		}

		// check if the field value can be changed
		// check if the field is exported
		if !field.CanSet() {
			continue
		}

		// if both have same kind then they can be placed together directly
		if field.Kind() == reflect.ValueOf(claimval).Kind() {
			switch field.Kind() {

			case reflect.Int:
				// convert value to int64 to check for overflow
				var i int64 = int64(claimval.(int))
				// checks if the int is overflowing in field type
				// if not then set value
				if !field.OverflowInt(i) {
					field.SetInt(i)
				}

			case reflect.Float32:
				// convert value to float64 to check for overflow
				var i float64 = claimval.(float64)
				//checks if the float will overflow in field type
				// if not then set value
				if !field.OverflowFloat(i) {
					field.SetFloat(i)
				}

			case reflect.Bool:
				var i bool = claimval.(bool)
				field.SetBool(i)

			case reflect.String:
				field.SetString(claimval.(string))

			case reflect.Slice, reflect.Array:
				v := reflect.ValueOf(claimval)
				// check if v is assignable to field type
				// if assignable then assign
				if v.Type().AssignableTo(field.Type()) {
					field.Set(v)
				}
			}
		}

		// if the field type requires struct type then recursively fill value
		if field.Kind() == reflect.Struct && reflect.ValueOf(claimval).Kind() == reflect.Map {
			_extractMapTo(claimval.(map[string]interface{}), field)
		}
	}
}
