package constant

import "time"

const (
	CLIENT_READ_LIMIT_SIZE = 4096
	CLIENT_READ_LIMIT_TIME = time.Minute * 10
)

const (
	USER_CHAT_EXPIRATION_DURATION = time.Hour * 24 * 30
)
