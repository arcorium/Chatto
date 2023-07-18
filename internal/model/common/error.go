package common

import "errors"

type Error struct {
	ErrorCode uint
	err       error
}

func (e Error) IsError() bool {
	return e.err != nil
}

func (e Error) Error() string {
	return e.err.Error()
}

// NewConditionalError Used to create common.Error based on error
func NewConditionalError(err error, badCode uint, badMsg string) Error {
	if err != nil {
		return NewError(badCode, badMsg)
	} else {
		return NoError()
	}
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
	INTERNAL_REPOSITORY_ERROR
	HASH_PASSWORD_ERROR
	CREATE_TOKEN_ERROR

	// User
	USER_NOT_FOUND_ERROR
	USER_SIGNUP_ERROR
	USER_ALREADY_ROOM_MEMBER
	USER_NOT_ROOM_MEMBER

	// Room
	ROOM_NOT_FOUND_ERROR
	ROOM_IS_PRIVATE_ERROR
	ROOM_CREATE_ERROR
	ROOM_NOT_EMPTY_ERROR

	// Auth
	AUTH_PASSWORD_INVALID_ERROR
	AUTH_TOKEN_NOT_FOUND_ERROR
	AUTH_TOKEN_BAD_OWNERSHIP_ERROR
	AUTH_TOKEN_NOT_VALIDATED_ERROR
	AUTH_TOKEN_UPDATE_ERROR
	AUTH_UNAUTHORIZED
)
