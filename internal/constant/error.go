package constant

// REST
// General
const (
	MSG_BAD_BODY_REQUEST      = "Some field in request body is missing"
	MSG_URI_PARAM_MISSING     = "URI parameter is missing"
	MSG_USER_AGENT_UNKNOWN    = "Unknown user-agent"
	MSG_INTERNAL_SERVER_ERROR = "Internal server error"
)

// User
const (
	MSG_FAILED_CREATE_USER = "Could not create new user"
	MSG_USER_NOT_FOUND     = "User doesn't exists"
	MSG_FAILED_UPDATE_USER = "Could not update user"
	MSG_FAILED_REMOVE_USER = "Could not remove user"
)

// Room
const (
	MSG_ROOM_CREATION_FAILED  = "Could not create new room"
	MSG_ROOM_NOT_FOUND        = "Room doesn't exist"
	MSG_REMOVE_NON_EMPTY_ROOM = "There are still users in the room"
	MSG_JOIN_PRIVATE_ROOM     = "Could not join into private room"
	MSG_JOIN_JOINED_ROOM      = "You are already room's member"
	MSG_NOT_ROOM_MEMBER       = "You are not room's member"
)

// Auth
const (
	MSG_FAILED_SIGNUP            = "Signup failed"
	MSG_FAILED_USER_LOGIN        = "Username or password does not match"
	MSG_LOGOUT_FAILED            = "Failed to logout"
	MSG_UNKNOWN_USER_AGENT       = "User agent is unknown"
	MSG_TOKEN_FIELD_INVALID_TYPE = "Token has invalid type field"
)

// Token
const (
	MSG_NO_ACCESS_TOKEN      = "Authentication required. Please provide a valid access token"
	MSG_TOKEN_NOT_FOUND      = "Token not found"
	MSG_BAD_FORMAT_TOKEN     = "Token have bad format. Please provide a valid token"
	MSG_NO_OWNER_TOKEN       = "Token is belong to no one"
	ERR_TOKEN_CREATION       = "Failed on token creation"
	MSG_TOKEN_REFRESH_FAILED = "Failed to refresh token"
	MSG_TOKEN_EXPIRED        = "Token you provided is expired"
)

// Chat
const (
	ERR_BAD_PAYLOAD      = "Payload is malformed"
	ERR_CLIENT_NOT_EXIST = "User is not exist"
	ERR_ROOM_NOT_EXIST   = "Room is not exist"
)
