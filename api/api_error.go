package api

import (
	"fmt"
)

// APIError ensures that the error message doesn't contain any sensitive
// inforation so that it can be safely displayed to a user.
type APIError struct {
	Message    string
	wrappedErr error
}

func NewAPIError(err error, msg string, args ...interface{}) APIError {
	return APIError{Message: fmt.Sprintf(msg, args...), wrappedErr: err}
}

func (a APIError) Error() string {
	return a.Message
}

func (a APIError) Unwrap() error {
	return a.wrappedErr
}
