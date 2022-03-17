package mongo

import "strings"

type Message struct {
	After             string             `json:"after"`
	Source            Source             `json:"source"`
	UpdateDescription *UpdateDescription `json:"updateDescription"`
	Op                string             `json:"op"`
	TsMS              int64              `json:"ts_ms"`
}

type Source struct {
	Version    string `json:"version"`
	Connector  string `json:"connector"`
	Name       string `json:"name"`
	TsMS       int64  `json:"ts_ms"`
	Snapshot   string `json:"snapshot"`
	DB         string `json:"db"`
	Rs         string `json:"rs"`
	Collection string `json:"collection"`
	Ord        int64  `json:"ord"`
}

type UpdateDescription struct {
	RemovedFields   interface{} `json:"removedFields"`
	TruncatedArrays interface{} `json:"truncatedArrays"`
	UpdatedFields   string      `json:"updatedFields"`
}

type (
	EventType string
)

var (
	Unknown EventType = "UNKNWON"
	Insert  EventType = "INSERT"
	Update  EventType = "UPDATE"
	Delete  EventType = "DELETE"
)

func (m *Message) EventType() EventType {
	switch m.Op {
	case "c":
		return Insert
	case "u":
		return Update
	case "d":
		return Delete
	default:
		return Unknown
	}
}

func (m *Message) TableName() string {
	if m.Source.Collection == "" {
		return m.Source.Collection
	}
	runes := []rune(m.Source.Collection)
	return strings.ToUpper(string(runes[0])) + string(runes[1:])
}
