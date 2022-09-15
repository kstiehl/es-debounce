package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/opensearch-project/opensearch-go"
)

var ErrOptNoAddress = errors.New("no address was specified")

// Client is the way to communicate with a configured opensearch.
type Client struct {
	opensearchClient *opensearch.Client
}

// NewWithDefaultClient creates a Client based on the given http.Client.
func NewWithDefaultClient() (Client, error) {
	client, err := opensearch.NewDefaultClient()
	if err != nil {
		return Client{}, err
	}
	return Client{
		opensearchClient: client,
	}, nil
}

// Document describes an indexable set of data.
type Document interface {
	// ID should return a unique ID for this document.
	ID() string

	// Index should return the target index
	Index() string

	// Data should return the docuemnt that should be indexed.
	Data() interface{}
}

// BulkIndex send multiple docuemnts to opensearch
func (client Client) BulkIndex(ctx context.Context, docs []Document) error {
	dataBytes, err := Bulk(docs).MarshalJSONToBuffer()
	if err != nil {
		return fmt.Errorf("unable to encode bulk %w", err)
	}

	// TODO: check response
	_, err = client.opensearchClient.Bulk(dataBytes)
	if err != nil {
		return fmt.Errorf("error during bulk index request to opensearch")
	}

	return nil
}

type Bulk []Document

func (b Bulk) MarshalJSONToBuffer() (*bytes.Buffer, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 512))
	for _, doc := range b {
		buffer.WriteString(`{"index": {"_index":"`)
		buffer.WriteString(doc.Index())
		buffer.WriteString(`", "_id": "`)
		buffer.WriteString(doc.ID())
		buffer.WriteString(`"}`)
		buffer.WriteRune('\n')

		encoder := json.NewEncoder(buffer)
		err := encoder.Encode(doc.Data())
		if err != nil {
			return nil, err
		}
	}

	return buffer, nil
}
