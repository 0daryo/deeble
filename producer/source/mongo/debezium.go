package mongo

import (
	"encoding/json"
	"strings"

	"github.com/0daryo/deeble/producer"
)

var _ producer.Producer = (*Producer)(nil)

type Producer struct{}

func (p *Producer) Produce(b []byte) ([]*producer.Message, error) {
	var m Message
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m.produce()
}

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

func (m *Message) eventType() producer.EventType {
	switch m.Op {
	case "c":
		return producer.Insert
	case "u":
		return producer.Update
	case "d":
		return producer.Delete
	default:
		return producer.Unknown
	}
}

func (m *Message) tableName() string {
	if m.Source.Collection == "" {
		return m.Source.Collection
	}
	runes := []rune(m.Source.Collection)
	return strings.ToUpper(string(runes[0])) + string(runes[1:])
}

func (m *Message) targets() (map[string]interface{}, error) {
	ret := map[string]interface{}{}
	if err := json.Unmarshal([]byte(m.After), &ret); err != nil {
		return nil, err
	}
	for k, v := range ret {
		// mongo bson can have start with $.
		if strings.HasPrefix(k, "$") {
			ret[k] = v
		}
	}
	return ret, nil
}

// mongo document can nest and bson has its type.
// e.g. map[_id:map[$oid:623bea8c0c02dba6bda13b63] first_name:hoge] to map[id:623bea8c0c02dba6bda13b63 first_name:hoge]
func parseNestedType(m map[string]interface{}) map[string]interface{} {
	msi := map[string]interface{}{}
	for k, v := range m {
		vv, ok := v.(map[string]interface{})
		if !ok {
			msi[k] = v
			continue
		}
		msi[k] = getFlat(vv, k)
	}
	return msi
}

// return nested first map value.
func getFlat(m map[string]interface{}, key string) interface{} {
	for _, v := range m {
		return v
	}
	return nil
}

func (m *Message) produce() ([]*producer.Message, error) {
	targets, err := m.targets()
	if err != nil {
		return nil, err
	}
	return []*producer.Message{
		{
			TableName: m.tableName(),
			EventType: m.eventType(),
			Targets:   parseNestedType(targets),
		},
	}, nil
}
