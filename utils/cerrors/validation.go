package cerrors

import (
	"encoding/json"
	"fmt"
)

type validationError struct {
	Code    string `json:"code"`
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   any    `json:"value"`
}

func NewValidation(code, field, message string, value any) *Error {
	return &Error{
		Value: validationError{Code: code, Field: field, Message: message, Value: value},
		next:  nil,
	}
}

func (e validationError) JSON() []byte {
	json, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return json
}

func (e validationError) CLI() string {
	return fmt.Sprint("Field `", e.Field, "` => ", e.Message)
}
