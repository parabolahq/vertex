package event

import "encoding/json"

type Event struct {
	AppName   string                  `json:"appName"`
	EventType string                  `json:"eventType"`
	Data      *map[string]interface{} `json:"data"`
}

func New(appName, eventType string, data any) Event {
	asMap := new(map[string]interface{})
	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &asMap)
	return Event{
		AppName:   appName,
		EventType: eventType,
		Data:      asMap,
	}
}
