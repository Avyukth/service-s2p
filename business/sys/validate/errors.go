package validate

import (
	"encoding/json"
	"errors"
)

var ErrInvalidID = errors.New("Id is not in its proper form")

type ErrorResponse struct {
	Error  string `json:"error"`
	Fields string `json:"fields,omitempty"`
}

type RequestError struct {
	Err    error
	Status int
	Fields error
}

func NewRequestError(err error, status int) error {
	return &RequestError{err, status, nil}
}

func (e *RequestError) Error() string {
	return e.Err.Error()
}

type FieldError struct {
	Field string `json:"field"`
	Erro  string `json:"error"`
}

type FieldErrors []FieldError

func (fe *FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}

func Cause(err error) error {
	root := err
	for {

		if err = errors.Unwrap(root); err == nil {
			return root
		}
		root = err
	}

}
