package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

var (
	ErrorEventIDEmpty        = errors.New("event IDs are not supposed to be empty")
	ErrorEventPayloadEmpty   = errors.New("event payload is nil")
	ErrorEventPayloadInvalid = errors.New("payload is not valid")
	ErrorNegativeStatusCode  = errors.New("opensearch replied with a negative status code")
)

const (
	logFieldStream  = "stream"
	logFieldEventID = "eventID"
)

type DataStream interface {
	Name() string
}

type EventPayload interface{}

type timestampedPayload struct {
	EventPayload
	TimeStamp time.Time `json:"@timestamp"`
}

func timestamp(payload EventPayload) timestampedPayload {
	return timestampedPayload{payload, time.Now().UTC()}
}

type Event struct {
	ID string

	Payload EventPayload
}

// IndexEvent takes a given Event and tries to appned it to the current stream
func IndexEvent(ctx context.Context, client Client, stream DataStream, event Event) error {
	log := logr.FromContextOrDiscard(ctx).
		WithName("opensearch-client").
		WithValues(logFieldStream, stream.Name(), logFieldEventID, event.ID)

	if err := validateEvent(event); err != nil {
		log.Info("event validation failed")
		return fmt.Errorf("failed to index document: %w", err)
	}

	timestamped := timestamp(event.Payload)
	payloadBytes, err := json.Marshal(timestamped)
	if err != nil {
		log.Info("converting payload to JSON failed")
		return fmt.Errorf("could not marshal payload to JSON: %w", err)
	}

	payloadReader := bytes.NewReader(payloadBytes)

	indexRequest := opensearchapi.CreateRequest{
		Index:      stream.Name(),
		DocumentID: event.ID,
		Body:       payloadReader,
	}

	response, err := indexRequest.Do(ctx, client)
	if err != nil {
		log.Info("executing request failed")
		return fmt.Errorf("executing index request failed: %w", err)
	}
	defer logClose(log, response.Body)

	if response.IsError() {
		analyzeBody(log, response)
		return ErrorNegativeStatusCode
	}
	return nil
}

// validateEvent checks whether the event can be safely processed.
func validateEvent(event Event) error {
	if event.ID == "" {
		return ErrorEventIDEmpty
	}

	if event.Payload == nil {
		return ErrorEventPayloadEmpty
	}

	return nil
}

// EnsureIndexTemplate makes sure that an Index Template is present and is configured in a given way.
//
// Note: If the configuration of an exisiting index template doesn't match the given configuration an error
// will be returned. Currently there is no save way for us to update the index template.
func EnsureIndexTemplate(ctx context.Context, client Client, config DataStream) error {
	log := logr.FromContextOrDiscard(ctx).WithName("opensearch-client").
		WithValues(logFieldStream, config.Name())

	streamName := config.Name()
	exists := opensearchapi.IndicesExistsIndexTemplateRequest{
		Name: streamName,
	}
	response, err := exists.Do(ctx, client)
	if err != nil {
		return err
	}
	defer logClose(log, response.Body)

	if response.StatusCode == http.StatusOK {
		log.Info("Datastream already present", "name", config.Name)
		return nil
	}

	bJson, err := json.Marshal([]string{streamName})
	if err != nil {
		log.Error(err, "unexpected error when marshalling index patterns slice to json")
		return err
	}
	indexTemplate := opensearchapi.IndicesPutIndexTemplateRequest{
		Body: strings.NewReader(`{
			"index_patterns": ` + string(bJson) + `,
			"data_stream": {},
			"priority": 100
		}
		`),
		Name: streamName,
	}

	response, err = indexTemplate.Do(ctx, client)
	if err != nil {
		log.Error(err, "error when executing request")
		return err
	}
	defer logClose(log, response.Body)

	if response.IsError() {
		analyzeBody(log, response)
		log.Info("unexpected status code", "statusCode", response.StatusCode)
		err = fmt.Errorf("unexpected response status code")
	}

	return nil
}

// analyzeBody dumps the reponse body to the log.
// sometimes opensearch replies with helpful error messages which can be useful when debugging.
func analyzeBody(log logr.Logger, response *opensearchapi.Response) {
	if response.IsError() {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Error(err, "failed reading opensearch response")
		}
		fields := append([]string{"payload", string(bodyBytes)})
		log.Info("elastic error resposne dump", "payload", fields)
	}
}

// logClose is a little helper to check the error when closing a response.
func logClose(log logr.Logger, closer io.Closer) {
	err := closer.Close()
	if err != nil {
		log.Error(err, "error when closing response")
	}
}
