package response

import "fmt"

var (
	ErrAuthFailed = ErrorResponse{
		Code:    "AUTH_FAILED",
		Message: "authentication failed",
	}
)

type ErrorResponse struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

func ErrNotFound(resource string) ErrorResponse {
	return ErrorResponse{
		Code:    "NOT_FOUND",
		Message: fmt.Sprintf("%s not found", resource),
	}
}

func ErrInvalid(invalid string) ErrorResponse {
	return ErrorResponse{
		Code:    "INVALID",
		Message: fmt.Sprintf("%s is invalid", invalid),
	}
}
