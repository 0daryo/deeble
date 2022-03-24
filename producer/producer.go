package producer

type (
	EventType string
)

var (
	Unknown EventType = "UNKNWON"
	Insert  EventType = "INSERT"
	Update  EventType = "UPDATE"
	Delete  EventType = "DELETE"
)

type Message struct {
	TableName string
	EventType EventType
	// key value pairs.
	// to insert, update, delete.
	Targets map[string]interface{}
}

func (m *Message) TargetKeys() []string {
	keys := make([]string, 0, len(m.Targets))
	for k := range m.Targets {
		keys = append(keys, k)
	}
	return keys
}

func (m *Message) TargetValues() []interface{} {
	vals := make([]interface{}, 0, len(m.Targets))
	for _, v := range m.Targets {
		vals = append(vals, v)
	}
	return vals
}

// producer may return multiple messages.
// e.g. 1 mongo document can be nest and split to multiple tables.
type Producer interface {
	Produce([]byte) ([]*Message, error)
}
