package util

import "chatto/internal/model"

func ForwardPayload[T any](senders []*model.Client, types string, outputData *T) {
	forwardPayload := model.NewPayloadOutput(types, outputData)
	for _, c := range senders {
		c.SendPayload(&forwardPayload)
	}
}

func SendErrorPayload(receiver *model.Client, code uint, message string) {
	payload := model.NewErrorPayloadOutput(code, message)
	receiver.SendPayload(&payload)
}

func SendSuccessPayload[T any](receiver *model.Client, outputData *T) {
	payload := model.NewPayloadOutput(model.PayloadSuccessResponse, outputData)
	receiver.SendPayload(&payload)
}
