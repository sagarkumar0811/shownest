package errors

// Code identifies the type of error.
type Code uint32

const (
	CodeUnknown Code = iota

	// General
	CodeInternal
	CodeUnimplemented
	CodeInvalidArgument
	CodeOutOfRange
	CodeFailedPrecondition
	CodeUnavailable
	CodeDeadlineExceeded
	CodeCanceled

	// Resource
	CodeNotFound
	CodeAlreadyExists
	CodeResourceExhausted

	// Auth
	CodeUnauthenticated
	CodePermissionDenied
	CodeTokenExpired
	CodeTokenInvalid
	CodeInvalidCredentials
	CodeUserBlocked

	// Validation
	CodeValidation

	// Database
	CodeDBError
	CodeDBNotFound
	CodeDBConflict
	CodeDBConnection
)

var codeStrings = map[Code]string{
	CodeUnknown:            "UNKNOWN",
	CodeInternal:           "INTERNAL",
	CodeUnimplemented:      "UNIMPLEMENTED",
	CodeInvalidArgument:    "INVALID_ARGUMENT",
	CodeOutOfRange:         "OUT_OF_RANGE",
	CodeFailedPrecondition: "FAILED_PRECONDITION",
	CodeUnavailable:        "UNAVAILABLE",
	CodeDeadlineExceeded:   "DEADLINE_EXCEEDED",
	CodeCanceled:           "CANCELED",
	CodeNotFound:           "NOT_FOUND",
	CodeAlreadyExists:      "ALREADY_EXISTS",
	CodeResourceExhausted:  "RESOURCE_EXHAUSTED",
	CodeUnauthenticated:    "UNAUTHENTICATED",
	CodePermissionDenied:   "PERMISSION_DENIED",
	CodeTokenExpired:       "TOKEN_EXPIRED",
	CodeTokenInvalid:       "TOKEN_INVALID",
	CodeInvalidCredentials: "INVALID_CREDENTIALS",
	CodeUserBlocked:        "USER_BLOCKED",
	CodeValidation:         "VALIDATION",
	CodeDBError:            "DB_ERROR",
	CodeDBNotFound:         "DB_NOT_FOUND",
	CodeDBConflict:         "DB_CONFLICT",
	CodeDBConnection:       "DB_CONNECTION",
}

func (c Code) String() string {
	if s, ok := codeStrings[c]; ok {
		return s
	}
	return "UNKNOWN"
}
