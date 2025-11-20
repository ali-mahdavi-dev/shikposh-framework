package phrases

type MessagePhrase string

const (
	// Default error IDs (Framework level)
	DefaultValidationID       MessagePhrase = "validation"
	DefaultNotFoundID         MessagePhrase = "not_found"
	DefaultConflictID         MessagePhrase = "conflict"
	DefaultUnauthorizedID     MessagePhrase = "unauthorized"
	DefaultForbiddenID        MessagePhrase = "forbidden"
	DefaultTimeoutID          MessagePhrase = "timeout"
	DefaultRateLimitID        MessagePhrase = "rate_limit"
	DefaultTooLargeID         MessagePhrase = "too_large"
	DefaultInternalID         MessagePhrase = "internal_server_error"
	DefaultMethodNotAllowedID MessagePhrase = "method_not_allowed"

	// Framework operation errors
	OperationCanNot MessagePhrase = "Operation.CanNot"

	// Framework parse errors
	FailedParseJson  MessagePhrase = "FailedParseJson"
	FailedParseQuery MessagePhrase = "FailedParseQuery"
	FailedParseForm  MessagePhrase = "FailedParseForm"
)
