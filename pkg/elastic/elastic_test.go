package elastic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestElasticSearch(t *testing.T) {
	t.Parallel()

	testingDoc := testingDoc{
		id:          "testingID",
		targetIndex: "testIndex",
		data: map[string]interface{}{
			"foo": "bar",
		},
	}

	t.Run("Index One", func(t *testing.T) {
		t.Parallel()

		b, err := Bulk([]Document{testingDoc}).MarshalJSONToBuffer()

		assert.NoError(t, err)
		assert.Equal(
			t,
			[]byte("{\"index\": {\"_index\":\"testIndex\", \"_id\": \"testingID\"}\n{\"foo\":\"bar\"}\n"),
			b.Bytes(),
		)
		fmt.Print(b)
	})

	t.Run("Index many", func(t *testing.T) {
		t.Parallel()

		secondDoc := testingDoc
		secondDoc.id = "second"

		b, err := Bulk([]Document{testingDoc, secondDoc}).MarshalJSONToBuffer()

		assert.NoError(t, err)
		assert.Equal(
			t,
			[]byte("{\"index\": {\"_index\":\"testIndex\", \"_id\": \"testingID\"}\n{\"foo\":\"bar\"}\n"+
				"{\"index\": {\"_index\":\"testIndex\", \"_id\": \"second\"}\n{\"foo\":\"bar\"}\n"),
			b.Bytes(),
		)
		fmt.Print(b)
	})
}

type testingDoc struct {
	targetIndex string
	id          string
	data        map[string]interface{}
}

func (t testingDoc) ID() string {
	return t.id
}

func (t testingDoc) Index() string {
	return t.targetIndex
}

func (t testingDoc) Data() interface{} {
	return t.data
}
