package services

import (
	"encoding/json"
	"strings"
)

//Error error interface
type Error interface {
	error
	Status() int
	ShouldDisplay() bool
	ErrorDetails() string
}

//APIError stores error information
type APIError struct {
	Message       string `json:"message"`
	MessageDetail string `json:"message-detail,omitempty"`
	Type          string `json:"type"`
	Code          int    `json:"-"`
	Private       bool   `json:"-"`
	InnerMessage  string `json:"-"`
}

//MarshalJSON custom marshaller for error information
func (e APIError) MarshalJSON() ([]byte, error) {
	type Alias APIError
	displayError := map[string]interface{}{
		"error": &struct {
			Alias
		}{
			Alias: (Alias)(e),
		},
	}

	return json.Marshal(displayError)
}

//NewError creates a new error object
func NewError(err error, message string, errorType string, private bool) *APIError {
	if err == nil {
		return nil
	}

	apiError := &APIError{
		Message: message,
		Type:    errorType,
		Private: private,
	}

	switch strings.ToUpper(errorType) {
	case "AUTHORIZATION":
		apiError.Code = 401
	case "DATABASECONNECTION":
		apiError.Code = 503
	case "DATABASEERROR":
		apiError.Code = 400
	case "FORMATERROR":
		apiError.Code = 400
	case "NOTFOUND":
		apiError.Code = 404
	case "VALIDATIONERROR":
		apiError.Code = 400
	default:
		apiError.Code = 500
	}

	apiError.InnerMessage = err.Error()

	if !private {
		apiError.MessageDetail = apiError.InnerMessage
	}

	return apiError
}

//Error returns error message
func (e APIError) Error() string {
	return e.Message
}

//ErrorDetails returns detailed error message
func (e APIError) ErrorDetails() string {
	return e.InnerMessage
}

//Status returns status code
func (e APIError) Status() int {
	return e.Code
}

//ShouldDisplay determines if error object should be displayed
func (e APIError) ShouldDisplay() bool {
	hideMessageCodes := []int{401, 404}

	for _, v := range hideMessageCodes {
		if v == e.Code {
			return false
		}
	}

	return true
}
