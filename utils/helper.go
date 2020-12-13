package utils

import (
	"encoding/json"
	"net/http"

	"oauth/constants"

	"go.mongodb.org/mongo-driver/bson"
)

// ErrorResponse struct defines the structure of the error response sent in API
type ErrorResponse struct {
	ErrorCode    constants.ErrorCode `json:"errorCode,omitempty"`
	ErrorMessage string              `json:"errorMessage,omitempty"`
}

type HTTPStatus struct {
	Code    int    `json:"statusCode"`
	Message string `json:"message"`
}

// respondJSON makes the response with payload as json format
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {

	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// respondError makes the error response with payload as json format
func RespondError(w http.ResponseWriter, httpResponseCode int, errorCode constants.ErrorCode, errorMessage string) {

	errorResponse := ErrorResponse{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}
	RespondJSON(w, httpResponseCode, errorResponse)
}

func GetTotalPagesCount(count, limit int) (pages int) {
	pages = count / limit // integer division, decimals are truncated
	if count%limit > 0 {
		pages = pages + 1
	}
	return
}

func ToBSONDoc(v interface{}) (doc *bson.M, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}
