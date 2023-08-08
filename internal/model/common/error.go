package common

import "errors"

type Error struct {
	ErrorCode uint
	err       error
}

func (e Error) IsError() bool {
	return e.err != nil
}

func (e Error) Error() error {
	return e.err
}

func (e Error) Message() string {
	return e.Error().Error()
}

// NewConditionalError Used to create common.Error based on error
func NewConditionalError(err error, badCode uint, badMsg string) Error {
	if err != nil {
		return NewError(badCode, badMsg)
	}
	return NoError()
}

func NewError(code uint, message string) Error {
	return Error{
		ErrorCode: code,
		err:       errors.New(message),
	}
}

func NoError() Error {
	return Error{
		ErrorCode: NO_ERROR,
		err:       nil,
	}
}

const (
	// General
	NO_ERROR = iota
	BAD_BODY_REQUEST_ERROR
	USER_AGENT_UNKNOWN_ERROR
	BAD_PARAMETER_ERROR
	INTERNAL_SERVER_ERROR
	HASH_PASSWORD_ERROR
	CREATE_TOKEN_ERROR

	// User
	USER_CREATION_ERROR
	USER_NOT_FOUND_ERROR
	USER_SIGNUP_ERROR
	USER_ALREADY_ROOM_MEMBER
	USER_NOT_ROOM_MEMBER
	MEMBER_ROOM_NOT_FOUND_ERROR

	// Room
	ROOM_NOT_FOUND_ERROR
	ROOM_IS_PRIVATE_ERROR
	ROOM_CREATE_ERROR
	ROOM_NOT_EMPTY_ERROR
	ROOM_ROLE_NOT_FOUND_ERROR
	ROOM_EMPTY_USER_ERROR

	// Auth
	AUTH_PASSWORD_INVALID_ERROR
	AUTH_TOKEN_NOT_FOUND_ERROR
	AUTH_TOKEN_BAD_OWNERSHIP_ERROR
	AUTH_TOKEN_NOT_VALIDATED_ERROR
	AUTH_TOKEN_UPDATE_ERROR
	AUTH_UNAUTHORIZED

	// Chat
	PAYLOAD_BAD_FORMAT_ERROR
)
