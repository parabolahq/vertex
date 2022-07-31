package error

import (
	"fmt"
	"tiny-submarine/event"
)

type ApiError struct {
	ErrorCode int    `json:"errorCode"`
	Data      string `json:"data"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("api error: (%d) %s", e.ErrorCode, e.Data)
}

func (e ApiError) AsEvent(appName string) event.Event {
	return event.New(appName, "error", e)
}

func (e ApiError) AsInternalEvent() event.Event {
	return e.AsEvent("tiny-submarine")
}
