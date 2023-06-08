package constant

import "time"

const (
	ACCESS_TOKEN_EXP_TIME  = time.Minute * 60 // TODO: Change this in production
	REFRESH_TOKEN_EXP_TIME = time.Hour * 24 * 90
	CONTEXT_TIMEOUT        = time.Second * 10
)
