package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

// validate holds the settings and caches for validating request struct values.
var validate *validator.Validate
var translator *ut.UniversalTranslator

func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}

	if err := validate.Struct(val); err != nil {

		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		// lang controls the language of the error messages. You could look at the
		// Accept-Language header if you intend to support multiple languages.

		var fields []FieldError
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
			}
			fields = append(fields, field)
		}

		errMessage := "field validation error"
		if len(fields) > 0 {
			errMessage = fields[0].Error
		}

		return &Error{
			Err:    errors.New(errMessage),
			Status: http.StatusBadRequest,
			Fields: fields,
		}
	}

	return nil
}

// Error is used to pass an error during the request through the
// application with web specific context.
type Error struct {
	Err    error
	Status int
	Fields []FieldError
}

func (e Error) Error() string {
	return e.Err.Error()
}

// NewRequestError wraps a provided error with an HTTP status code. This
// function should be used when handlers encounter expected errors.
func NewRequestError(err error, status int) error {
	return &Error{err, status, nil}
}

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field,omitempty"`
	Error string `json:"error"`
}

// ErrorResponse is the form used for API responses from failures in the API.
type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}
