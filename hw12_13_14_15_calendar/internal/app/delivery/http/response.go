package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator"
)

type (
	response struct {
		Data   responseData   `json:"data"`
		Errors responseErrors `json:"errors"`
	}
	responseData   interface{}
	responseErrors map[string]string
)

func makeResponseError(w http.ResponseWriter, r *http.Request, code int, err error) {
	respond(w, code, nil, err)
}

func makeResponse(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	respond(w, code, data, nil)
}

func respond(w http.ResponseWriter, code int, data interface{}, err error) {
	w.WriteHeader(code)

	response := &response{
		Data:   data,
		Errors: createErrors(err),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
	}
}

func createErrors(err error) responseErrors {
	if err == nil {
		return *new(responseErrors)
	}

	validationError := new(validator.ValidationErrors)

	if errors.As(err, validationError) {
		errors := make(responseErrors)
		for _, e := range *validationError {
			errors[e.Field()] = e.Translate(nil)
		}

		return errors
	}

	return responseErrors{
		"error": err.Error(),
	}
}
