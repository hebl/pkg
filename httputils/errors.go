package httputils

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// GetHTTPErrorStatusCode retrieves status code from error message.
func GetHTTPErrorStatusCode(err error) int {
	if err == nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("unexpected HTTP error handling")
		return http.StatusInternalServerError
	}

	//var statusCode int

	if e, ok := err.(StatusCoder); ok {
		return e.StatusCode()
	}
	return http.StatusInternalServerError
}

// MakeErrorHandler makes an HTTP handler that decodes a Docker error and
// returns it in the response.
func MakeErrorHandler(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode := GetHTTPErrorStatusCode(err)
		response := &ErrorResponse{
			Code:    statusCode,
			Message: err.Error(),
		}
		WriteJSON(w, statusCode, response)

	}
}

//StatusCoder error response
type StatusCoder interface {
	StatusCode() int
}

// ErrorResponse error response
type ErrorResponse struct {
	Code    int    `json:"statusCode"`
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

//StatusCode statuscode
func (e ErrorResponse) StatusCode() int {
	return e.Code
}
