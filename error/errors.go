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
	return e.AsEvent("submarine")
}
