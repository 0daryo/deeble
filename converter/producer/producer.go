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
	Updates   map[string]interface{}
	Inserts   map[string]interface{}
	DeleteKey interface{}
}
