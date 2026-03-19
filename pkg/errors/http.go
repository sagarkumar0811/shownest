package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

var codeHTTPStatus = map[Code]int{
	CodeUnknown:            http.StatusInternalServerError,
	CodeInternal:           http.StatusInternalServerError,
	CodeUnimplemented:      http.StatusNotImplemented,
	CodeInvalidArgument:    http.StatusBadRequest,
	CodeOutOfRange:         http.StatusBadRequest,
	CodeFailedPrecondition: http.StatusBadRequest,
	CodeUnavailable:        http.StatusServiceUnavailable,
	CodeDeadlineExceeded:   http.StatusGatewayTimeout,
	CodeCanceled:           http.StatusBadRequest,
	CodeNotFound:           http.StatusNotFound,
	CodeAlreadyExists:      http.StatusConflict,
	CodeResourceExhausted:  http.StatusTooManyRequests,
	CodeUnauthenticated:    http.StatusUnauthorized,
	CodePermissionDenied:   http.StatusForbidden,
	CodeTokenExpired:       http.StatusUnauthorized,
	CodeTokenInvalid:       http.StatusUnauthorized,
	CodeInvalidCredentials: http.StatusUnauthorized,
	CodeUserBlocked:        http.StatusForbidden,
	CodeValidation:         http.StatusUnprocessableEntity,
	CodeDBError:            http.StatusInternalServerError,
	CodeDBNotFound:         http.StatusNotFound,
	CodeDBConflict:         http.StatusConflict,
	CodeDBConnection:       http.StatusServiceUnavailable,
}

// HTTPStatus returns the HTTP status code for the given error code.
func HTTPStatus(code Code) int {
	if status, ok := codeHTTPStatus[code]; ok {
		return status
	}
	return http.StatusInternalServerError
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorBody `json:"error"`
}

// WriteError writes a JSON error response derived from err.
// If err is an *AppError, its code and message are used; otherwise a generic internal error is written.
func WriteError(w http.ResponseWriter, err error) {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		appErr = &AppError{Code: CodeInternal, Message: "internal server error"}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(HTTPStatus(appErr.Code))
	json.NewEncoder(w).Encode(errorResponse{ //nolint:errcheck
		Error: errorBody{
			Code:    appErr.Code.String(),
			Message: appErr.Message,
		},
	})
}
