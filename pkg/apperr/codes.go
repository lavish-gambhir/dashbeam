package apperr

type ErrCode string

const (
	// Common Error Codes
	NoErr              ErrCode = "NO_ERROR"
	Unknown            ErrCode = "UNKNOWN"
	Internal           ErrCode = "INTERNAL"
	BadRequest         ErrCode = "BAD_REQUEST"
	NotFound           ErrCode = "NOT_FOUND"
	Unauthorized       ErrCode = "UNAUTHORIZED"
	Forbidden          ErrCode = "FORBIDDEN"
	Conflict           ErrCode = "CONFLICT"
	ServiceUnavailable ErrCode = "SERVICE_UNAVAILABLE"

	// Authentication Specific Error Codes
	InvalidCredentials ErrCode = "AUTH_INVALID_CREDENTIALS"
	TokenExpired       ErrCode = "AUTH_TOKEN_EXPIRED"
	InvalidToken       ErrCode = "AUTH_INVALID_TOKEN"
	UserNotFound       ErrCode = "AUTH_USER_NOT_FOUND"
	UserAlreadyExists  ErrCode = "AUTH_USER_ALREADY_EXISTS"

	// Database Specific Error Codes
	DBConnectionFailed ErrCode = "DB_CONNECTION_FAILED"
	DBQueryFailed      ErrCode = "DB_QUERY_FAILED"
	DBRecordNotFound   ErrCode = "DB_RECORD_NOT_FOUND"
	DBDuplicateEntry   ErrCode = "DB_DUPLICATE_ENTRY"

	// Validation Specific Error Codes
	ValidationFailed     ErrCode = "VALIDATION_FAILED"
	InvalidFormat        ErrCode = "VALIDATION_INVALID_FORMAT"
	MissingRequiredField ErrCode = "VALIDATION_MISSING_REQUIRED_FIELD"
	JSONEncodingFailed   ErrCode = "JSON_ENCODING_FAILED"
	JSONDecodingFailed   ErrCode = "JSON_DECODING_FAILED"

	// Redis Error Codes
	RedisPipeExecFailed ErrCode = "REDIS_PIPE_EXEC_FAILED"
	RedisNoEvent        ErrCode = "REDIS_NO_EVENT"
	RedisUnknown        ErrCode = "UNKNOWN_REDIS_ERR"
)
