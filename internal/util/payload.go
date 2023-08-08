package util

import (
	"chatto/internal/model"
	"chatto/internal/model/common"
)

func SendErrorPayload(receiver *model.Client, err common.Error) {
	payload := model.NewErrorPayloadOutput(err.ErrorCode, err.Message())
	receiver.SendPayload(&payload)
}

func SendSuccessPayload[T any](receiver *model.Client, outputData *T) {
	payload := model.NewPayloadOutput(model.PayloadSuccessResponse, outputData)
	receiver.SendPayload(&payload)
}

func SendNilSuccessPayload(receiver *model.Client) {
	var t *int
	payload := model.NewPayloadOutput(model.PayloadSuccessResponse, t)
	receiver.SendPayload(&payload)
}
