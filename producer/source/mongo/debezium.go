package mongo

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/0daryo/deeble/producer"
	"github.com/google/uuid"
)

var _ producer.Producer = (*Producer)(nil)

type Producer struct {
}

func (p *Producer) Produce(b []byte) ([]*producer.Message, error) {
	var m message
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m.produce()
}

type (
	message struct {
		After             string             `json:"after"`
		Source            source             `json:"source"`
		UpdateDescription *updateDescription `json:"updateDescription"`
		Op                string             `json:"op"`
		TsMS              int64              `json:"ts_ms"`
	}
	messageConverter struct {
		tableName string
		eventType producer.EventType
		targets   map[string]interface{}
		depth     int64
	}
	messageConverters []*messageConverter
)

func (mcs *messageConverters) slice() []*messageConverter {
	if mcs == nil {
		return nil
	}
	return []*messageConverter(*mcs)
}

type source struct {
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

type updateDescription struct {
	RemovedFields   interface{} `json:"removedFields"`
	TruncatedArrays interface{} `json:"truncatedArrays"`
	UpdatedFields   string      `json:"updatedFields"`
}

func eventType(op string) producer.EventType {
	switch op {
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

func tableName(collectionName string) string {
	if collectionName == "" {
		return ""
	}
	runes := []rune(collectionName)
	return strings.ToUpper(string(runes[0])) + string(runes[1:])
}

// mongo document can nest and bson has its type.
// e.g. map[_id:map[$oid:623bea8c0c02dba6bda13b63] first_name:hoge] to map[id:623bea8c0c02dba6bda13b63 first_name:hoge]
func parseNestedType(m map[string]interface{}) map[string]interface{} {
	msi := map[string]interface{}{}
	for k, v := range m {
		vv, ok := v.(map[string]interface{})
		if !ok || !hasSpecialCharKey(vv) {
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

// if map key has bson special characters
func hasSpecialCharKey(m map[string]interface{}) bool {
	for k := range m {
		if strings.HasPrefix(k, "$") {
			return true
		}
	}
	return false
}

func (m *message) produce() ([]*producer.Message, error) {
	after := make(map[string]interface{})
	if err := json.Unmarshal([]byte(m.After), &after); err != nil {
		return nil, err
	}
	targets := after
	mcs := messageConverters{}
	id := id(targets)
	delete(targets, "_id")
	mcs.fill(targets, m.Source.Collection, 0, eventType(m.Op), map[string]interface{}{
		"Id": id,
	})
	mcsSlice := mcs.slice()
	sort.SliceStable(mcsSlice, func(i, j int) bool { return mcsSlice[i].depth < mcsSlice[j].depth })
	return messageConverters(mcsSlice).producerMessge(), nil
}

func id(m map[string]interface{}) interface{} {
	mm := parseNestedType(m)
	id, ok := mm["_id"]
	if !ok {
		return genID()
	}
	switch id.(type) {
	case map[string]interface{}:
		if oid, ok := id.(map[string]interface{})["$oid"]; ok {
			if str, ook := oid.(string); ook {
				return str
			}
			return genID()
		}
		return genID()
	default:
		return id
	}
}

// fill build messageConverters with provided targets recursively.
func (mcs *messageConverters) fill(m map[string]interface{}, collectionName string, depth int64, eventType producer.EventType, fks map[string]interface{}) {
	if mcs == nil {
		return
	}
	mm := parseNestedType(m)
	// add foreign keys for interleave.
	for k, v := range fks {
		mm[k] = v
	}
	if nonest(mm) {
		*mcs = append(*mcs, &messageConverter{
			tableName: tableName(collectionName),
			eventType: eventType,
			targets:   mm,
			depth:     depth,
		})
		return
	}
	mc := &messageConverter{
		tableName: tableName(collectionName),
		eventType: eventType,
		depth:     depth,
	}
	targets := make(map[string]interface{})
	for k, v := range mm {
		switch v.(type) {
		case map[string]interface{}:
			nextFK := nextFK(tableName(collectionName), fks)
			mcs.fill(v.(map[string]interface{}), tableName(k), depth+1, eventType, nextFK)
		default:
			targets[k] = v
		}
	}
	mc.targets = targets
	*mcs = append(*mcs, mc)
}

func nextFK(tableName string, fks map[string]interface{}) map[string]interface{} {
	nextFK := make(map[string]interface{})
	for k, v := range fks {
		if k == "Id" {
			nextFK[parentFK(tableName)] = v
			continue
		}
		nextFK[k] = v
	}
	nextFK["Id"] = genID()
	return nextFK
}

func parentFK(tableName string) string {
	singlar := strings.TrimSuffix(tableName, "s")
	return singlar + "Id"
}

// check if map has no nested map[string]interface{}
func nonest(m map[string]interface{}) bool {
	for _, v := range m {
		switch v.(type) {
		case map[string]interface{}:
			return false
		}
	}
	return true
}

// convert messageConverters to producer.Messages slice
func (mcs messageConverters) producerMessge() []*producer.Message {
	ms := []*producer.Message{}
	for _, mc := range mcs {
		ms = append(ms, &producer.Message{
			TableName: mc.tableName,
			EventType: mc.eventType,
			Targets:   mc.targets,
		})
	}
	return ms
}

var genID = func() string {
	uuid, _ := uuid.NewRandom()
	return uuid.String()
}
