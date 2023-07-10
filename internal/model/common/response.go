package common

func NewErrorResponse(code uint, message string, data any) Response {
	return Response{
		Status:  "error",
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewSuccessResponse(code uint, data any) Response {
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
