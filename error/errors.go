package error

import (
	"fmt"
	"vertex/communication"
)

type ApiError struct {
	ErrorCode int    `json:"errorCode"`
	Data      string `json:"data"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("api error: (%d) %s", e.ErrorCode, e.Data)
}

func (e ApiError) AsEvent(appName string) communication.Event {
	return communication.New(appName, "error", e)
}

func (e ApiError) AsInternalEvent() communication.Event {
	return e.AsEvent("vertex")
}

func InternalError() ApiError {
	return ApiError{
		ErrorCode: 0,
		Data:      "Internal error occurred",
	}
}

func InvalidToken() ApiError {
	return ApiError{
		ErrorCode: 1,
		Data:      "Invalid or expired token",
	}
}

func BadRequest() ApiError {
	return ApiError{
		ErrorCode: 2,
		Data:      "Parse of JSON body failed",
	}
}

func SendingMQError() ApiError {
	return ApiError{
		ErrorCode: 3,
		Data:      "Failed to send message to RabbitMQ",
	}
}
