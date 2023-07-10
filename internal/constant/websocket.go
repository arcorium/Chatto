package constant

import "time"

const (
	CLIENT_READ_LIMIT_SIZE = 4096
	CLIENT_READ_LIMIT_TIME = time.Minute * 10
)

// Error Code
const (
	NO_RECEIVER = iota
	INTERNAL_ERROR
)

// Error Message
const (
	NO_RECEIVER_MSG = "No receiver"
)
