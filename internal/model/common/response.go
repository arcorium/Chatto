package common

func NewErrorResponse(err Error) Response {
	return Response{
		Status:  "error",
		Code:    err.ErrorCode,
		Message: err.err.Error(),
		Data:    nil,
	}
}

func NewSuccessResponse(err Error, data any) Response {
	return Response{
		Status: "success",
		Code:   err.ErrorCode,
		Data:   data,
	}
}

func NewSuccessResponseWithCode(code uint, data any) Response {
	return Response{
		Status: "success",
		Code:   code,
		Data:   data,
	}
}

type Response struct {
	Status  string `json:"status"`
	Code    uint   `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}
