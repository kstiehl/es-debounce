package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-logr/logr"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

// DataStreamConfig holds all information in order to create a data stream.
type DataStreamConfig struct {
	// Name holds the name of the data stream.
	Name string

	// IndexPatterns holds the patterns that should be defined for the data stream
	IndexPatterns []string
}

// EnsureIndexTemplate makes sure that an Index Template is present and is configured in a given way.
//
// Note: If the configuration of an exisiting index template doesn't match the given configuration an error
// will be returned. Currently there is no save way for us to update the index template.
func EnsureIndexTemplate(ctx context.Context, client *opensearch.Client, config DataStreamConfig) error {
	log := logr.FromContextOrDiscard(ctx).WithName("opensearch-client")
	exists := opensearchapi.IndicesExistsIndexTemplateRequest{
		Name: config.Name,
	}
	res, err := exists.Do(ctx, client)
	if err != nil {
		return err
	}
	defer logClose(log, res.Body)

	if res.StatusCode == http.StatusOK {
		log.Info("Datastream already present", "name", config.Name)
		return nil
	}

	bJson, err := json.Marshal(&config.IndexPatterns)
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
		Name: config.Name,
	}

	res, err = indexTemplate.Do(ctx, client)
	if err != nil {
		log.Error(err, "error when executing request")
		return err
	}
	defer logClose(log, res.Body)

	if res.IsError() {
		bBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Error(err, "unable to read response body")
		}

		err = fmt.Errorf("unexpected response status code ")
		log.Error(err, "unexpected response when creating index", "responseBody", string(bBody))
	}

	return nil
}

// logClose is a little helper to check the error when closing a response.
func logClose(log logr.Logger, closer io.Closer) {
	err := closer.Close()
	if err != nil {
		log.Error(err, "error when closing response")
	}
}
