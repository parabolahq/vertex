package event

import "encoding/json"

type Event struct {
	ServiceAlias string                  `json:"appName"`
	EventType    string                  `json:"eventType"`
	Data         *map[string]interface{} `json:"data"`
}

func New(serviceAlias, eventType string, data any) Event {
	asMap := new(map[string]interface{})
	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &asMap)
	return Event{
		ServiceAlias: serviceAlias,
		EventType:    eventType,
		Data:         asMap,
	}
}
