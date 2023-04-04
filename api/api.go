package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/go-logr/logr"
	"github.com/kstiehl/index-bouncer/grpc/types"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

const (
	// TargetIndexName holds the name to which all events are currently written.
	// This is hardcoded for now.
	TargetIndexName = "eventingest"
)

// Prepare should prepare the underlying storage system in order that this stream can be used.
func Index(ctx context.Context, client *opensearch.Client, event *types.Event) error {
	log := logr.FromContextOrDiscard(ctx).WithName("API")

	reader, err := serializeEvent(event)
	if err != nil {
		return err
	}

	indexRequest := opensearchapi.CreateRequest{
		Index:      "eventingest",
		DocumentID: event.EventID,
		Body:       reader,
	}

	response, err := indexRequest.Do(ctx, client)
	if err != nil {
		return err
	}

	logClose(log, response.Body)

	if response.IsError() {

		analyzeBody(log, response)
		log.Info("unexpected response status code from opensearch",
			"statusCode", response.StatusCode)

		// error message from opensearch should never be leaked to client
		// note: maybe it makes sense to include EventID in the future.
		// Could be helpful when debugging.
		return errors.New("Failed to index event")
	}

	return nil
}

// serializeEvent should convert an Event to JSON.
// Since this function will invoked quite alot it won't hurt
// to pay attention to memory consumption and performance
func serializeEvent(event *types.Event) (io.Reader, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 4096)) // could be a place for sync.Pool{} but needs more benchmarks
	buffer.WriteString("{\"eventID\": \"")           // 13
	buffer.WriteString(event.EventID)
	buffer.WriteString("\", \"objectID\": \"") // 15
	buffer.WriteString(event.ObjectID)
	buffer.WriteString("\", \"data\": [") // 2
	for i, value := range event.Data {
		if i != 0 {
			buffer.WriteRune(',')
		}

		buffer.WriteString("{\"")
		buffer.WriteString(value.Key)
		buffer.WriteString("\": ")
		err := json.NewEncoder(buffer).Encode(getEventDataValue(value))
		if err != nil {
			return nil, fmt.Errorf("unable to serialize event: %w", err)
		}
		buffer.WriteString("}")
	}

	buffer.WriteString("]}")
	return buffer, nil
}

func getEventDataValue(v *types.EventData) interface{} {
	switch value := v.Value.(type) {
	case *types.EventData_StringValue:
		return value.StringValue
	case *types.EventData_BoolValue:
		return value.BoolValue
	case *types.EventData_NumberValue:
		return value.NumberValue
	}
	return ""
}

// logClose is a little helper to check the error when closing a response.
func logClose(log logr.Logger, closer io.Closer) {
	err := closer.Close()
	if err != nil {
		log.Error(err, "error when closing response")
	}
}

// analyzeBody dumps the reponse body to the log.
// sometimes opensearch replies with helpful error messages which can be useful when debugging.
func analyzeBody(log logr.Logger, response *opensearchapi.Response) {
	if response.IsError() {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Error(err, "failed reading opensearch response")
		}
		fields := append([]string{"payload"}, string(bodyBytes))
		log.Info("elastic error resposne dump", "payload", fields)
	}
}
