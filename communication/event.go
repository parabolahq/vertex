package communication

import (
	"encoding/json"
)

// Event is data, that is sent to user by pool
type Event struct {
	ServiceAlias string                  `json:"serviceAlias"`
	EventType    string                  `json:"eventType"`
	RecipientIds []string                `json:"recipientIds,omitempty"`
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

func (e *Event) AsBytes() (data []byte) {
	data, _ = json.Marshal(e)
	return
}
