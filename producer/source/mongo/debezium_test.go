package mongo

import (
	"os"
	"reflect"
	"testing"

	"github.com/0daryo/deeble/producer"
	"github.com/stretchr/testify/assert"
)

func TestProducer_Produce(t *testing.T) {
	testInsert, err := os.ReadFile("testdata/insert.json")
	assert.Nil(t, err)
	testUpdate, err := os.ReadFile("testdata/update.json")
	assert.Nil(t, err)
	testBson, err := os.ReadFile("testdata/objectid.json")
	assert.Nil(t, err)
	testNest, err := os.ReadFile("testdata/nest.json")
	assert.Nil(t, err)
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		p       *Producer
		args    args
		want    []*producer.Message
		wantErr bool
	}{
		{
			name: "success insert",
			p:    &Producer{},
			args: args{
				b: testInsert,
			},
			want: []*producer.Message{
				{
					TableName: "Customers",
					EventType: producer.Insert,
					Targets: map[string]interface{}{
						"_id":       float64(1015),
						"firstName": "hoge",
						"lastName":  "fuga",
						"email":     "hoge@example.com",
					},
				},
			},
		},
		{
			name: "success update",
			p:    &Producer{},
			args: args{
				b: testUpdate,
			},
			want: []*producer.Message{
				{
					TableName: "Customers",
					EventType: producer.Update,
					Targets: map[string]interface{}{
						"_id":       float64(1015),
						"firstName": "hoge",
						"lastName":  "fuga",
						"email":     "hoge@example.com",
					},
				},
			},
		},
		{
			name: "success bson",
			p:    &Producer{},
			args: args{
				b: testBson,
			},
			want: []*producer.Message{
				{
					TableName: "Customers",
					EventType: producer.Update,
					Targets: map[string]interface{}{
						"_id":        "623bea8c0c02dba6bda13b63",
						"first_name": "hoge",
					},
				},
			},
		},
		{
			name: "success nested",
			p:    &Producer{},
			args: args{
				b: testNest,
			},
			want: []*producer.Message{
				{
					TableName: "Customers",
					EventType: producer.Update,
					Targets: map[string]interface{}{
						"_id":        "623d8883f25162b8f356ce91",
						"first_name": "mike",
					},
				},
				{
					TableName: "Nest",
					EventType: producer.Update,
					Targets: map[string]interface{}{
						"last_name": "fuga",
					},
				},
				{
					TableName: "Nest1",
					EventType: producer.Update,
					Targets: map[string]interface{}{
						"hoge": "fuga",
					},
				},
			},
		},
		{
			name: "invalid json",
			p:    &Producer{},
			args: args{
				b: []byte("invalid json"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Producer{}
			got, err := p.Produce(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Producer.Produce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Producer.Produce() = %v, want %v", got, tt.want)
			}
		})
	}
}
