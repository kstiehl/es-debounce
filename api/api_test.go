package api

import (
	"testing"

	. "github.com/kstiehl/index-bouncer/grpc/types"
)

func BenchmarkSerialization(b *testing.B) {
	event := Event{
		EventID:  "testrelkglrtekly",
		ObjectID: "dskjggjktrjhrt",
		Data: []*EventData{
			{Key: "eventData1.com.io", Value: &EventData_StringValue{"dksfgkrnegkret"}},
			{Key: "ejfkrjge.edor", Value: &EventData_NumberValue{5959}},
			{Key: "ejfkejrekjk.frogrejgjt", Value: &EventData_BoolValue{true}},
		},
	}

	for i := 0; i < b.N; i++ {
		serializeEvent(&event)
	}
}
