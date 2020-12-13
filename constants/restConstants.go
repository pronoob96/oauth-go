package constants

// ErrorCode defines the codes for different errors in billing service
type ErrorCode int

// General errors
const (
	IncompleteTokenError ErrorCode = 6001
)

const (
	Success = "success"
	AccessTokenHeaderName = "x-access-token"
)
