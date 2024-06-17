package dto

import (
	"encoding/json"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

// DecodeValid decodes the request body into a value of type T and returns an error if the decoding fails.
func DecodeValid[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, errors.Wrap(err, "decode json failed")
	}
	return v, nil
}

// Encode encodes the value v into JSON format and writes it to the http.ResponseWriter.
// It sets the Content-Type header to "application/json" and the status code to the given status.
// If the encoding fails, it returns an error wrapped with a message "encode json failed".
func Encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return errors.Wrap(err, "encode json failed")
	}
	return nil
}

// SendErrorJsonResponse sends a JSON error response with the given error message
// or map of error messages to the http.ResponseWriter. It uses the Encode function to encode the error
func SendErrorJsonResponse[T string | map[string]string](w http.ResponseWriter, logger *zap.Logger, resErr T) {
	errObj := ErrorResponse[T]{
		Error: resErr,
	}
	err := Encode[ErrorResponse[T]](w, http.StatusBadRequest, errObj)
	if err != nil {
		logger.Error("could not encode error response",
			zap.Error(err))
	}
}
