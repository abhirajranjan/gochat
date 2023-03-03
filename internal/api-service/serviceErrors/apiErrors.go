package serviceErrors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type IErr interface {
	Error() string
	To_json() ([]byte, error)
}

type BindingErr struct {
	Params     string `json:"params"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *BindingErr) To_json() ([]byte, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func (e *BindingErr) Error() string {
	return fmt.Sprintf("%s: %s", e.Params, e.Message)
}

func NewBindingErr(params string, message string) IErr {
	return &BindingErr{Params: params, Message: message, StatusCode: http.StatusBadRequest}
}

///////////////////////////////////////////////////////////////////////////////////////////

type ErrorArray []IErr

var rwmutex = sync.RWMutex{}

func (e *ErrorArray) To_json() ([]byte, error) {
	e.getWriteLock()
	defer e.releaseWriteLock()

	_def_e := *e
	arr := []map[string]interface{}{}

	for _, ierr := range _def_e {
		if ierr == nil {
			continue
		}
		temp := map[string]interface{}{}
		s, err := ierr.To_json()

		if err != nil {
			return []byte{}, err
		}
		if err := json.Unmarshal(s, &temp); err != nil {
			return []byte{}, err
		}
		arr = append(arr, temp)
	}

	return json.Marshal(arr)
}

func (e *ErrorArray) Error() string {
	err := ""
	for i, error := range *e {
		err += fmt.Sprintf("%d: %s\n", i, error.Error())
	}
	return err
}

func (e *ErrorArray) getWriteLock() {
	rwmutex.Lock()
}

func (e *ErrorArray) releaseWriteLock() {
	rwmutex.Unlock()
}

///////////////////////////////////////////////////////////////////////////////////////////

type ValidationErr struct {
	Param      string `json:"param"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *ValidationErr) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.StatusCode, e.Param, e.Message)
}

func (e *ValidationErr) To_json() ([]byte, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func NewValidationErr(param string, message string) IErr {
	return &ValidationErr{Param: param, Message: message, StatusCode: http.StatusUnprocessableEntity}
}

////////////////////////////////////////////////////////////////////////////////////////////

type StandardErr struct {
	Param             string        `json:"param"`
	Message           string        `json:"message"`
	AdditionalDetails []interface{} `json:"additionalDetails"`
}

func (e *StandardErr) Error() string {
	return fmt.Sprintf("%s:%s", e.Param, e.Message)
}

func (e *StandardErr) To_json() ([]byte, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func NewStandardErr(param string, message string, additionalDetails ...interface{}) IErr {
	return &StandardErr{Param: param, Message: message, AdditionalDetails: additionalDetails}
}
