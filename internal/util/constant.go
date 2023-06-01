package util

import "time"

const (
	ACCESS_TOKEN_EXP_TIME  = time.Minute * 60 // TODO: Change this in production
	REFRESH_TOKEN_EXP_TIME = time.Hour * 24 * 90
	CONTEXT_TIMEOUT        = time.Second * 10
)

const (
	CLIENT_READ_LIMIT_SIZE = 4096
	CLIENT_READ_LIMIT_TIME = time.Minute * 10
)

const (
	ERR_BODY_REQUEST_MISSING  = "Some field in request body is missing"
	ERR_URI_PARAM_MISSING     = "URI parameter is missing"
	ERR_INVALID_TYPE          = "Invalid type"
	ERR_USER_AGENT_MIDDLEWARE = "Unknown user-agent"
)

const (
	ERR_USER_CREATE    = "Could not create new user"
	ERR_USER_NOT_FOUND = "User doesn't exists"
	ERR_USER_UPDATE    = "Could not update user"
	ERR_USER_REMOVE    = "Could not remove user"
)

const (
	ERR_SIGNUP                  = "Signup failed"
	ERR_LOGIN_USERNAME_PASSWORD = "Username or password does not match"
	ERR_LOGOUT                  = "Failed to logout"
)

const (
	ERR_NO_ACCESS_TOKEN = "Authentication required. Please provide a valid access token"
	ERR_TOKEN_FORMAT    = "Token have bad format. Please provide a valid token"
	ERR_TOKEN_NO_OWNER  = "Token is belong to no one"
	ERR_TOKEN_FIELD     = "Token have empty required field"
	ERR_TOKEN_CREATION  = "Failed on token creation"
	ERR_TOKEN_REFRESH   = "Failed to refresh token"
)
