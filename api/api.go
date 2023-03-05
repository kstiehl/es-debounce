package api

import (
	"context"
	"fmt"
	"io"

	"github.com/go-logr/logr"
	"github.com/kstiehl/index-bouncer/grpc"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

// Prepare should prepare the underlying storage system in order that this stream can be used.
func Prepare(ctx context.Context, client *opensearch.Client, event *grpc.Event) error {
	log := logr.FromContextOrDiscard(ctx).WithName("API")
	indexRequest := opensearchapi.CreateRequest{
		Index:      "eventingest",
		DocumentID: event.EventID,
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
		err = fmt.Errorf("Failed to index event")
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
