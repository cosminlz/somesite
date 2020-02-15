package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

// GenericError represents structure for generic errors. All responses should be the same
type GenericError struct {
	Code  int         `json:"code"`
	Error string      `json:"error"`
	Data  interface{} `json:"data,omitempty"`
}

// WriteError writes a error response
func WriteError(w http.ResponseWriter, code int, message string, data interface{}) {
	response := GenericError{
		Code:  code,
		Error: message,
		Data:  data,
	}

	WriteJSON(w, code, response)
}

// WriteJSON used to write any generic json object to response
func WriteJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logrus.WithError(err).Warn("Error writing json response")
	}
}
