package api

import (
	"context"
	"fmt"

	"github.com/kstiehl/index-bouncer/pkg/opensearch"
)

type (
	Stream = opensearch.DataStream
	Client = opensearch.Client
	Event  = opensearch.Event
)

type API struct {
	client Client
}

// NewAPI returns an API object with the given configuration
func NewAPI(client Client) API {
	return API{client}
}

// Prepare should prepare the underlying storage system in order that this stream can be used.
func (api API) Prepare(ctx context.Context, stream Stream) error {
	err := opensearch.EnsureIndexTemplate(ctx, api.client, stream)
	if err != nil {
		return fmt.Errorf("error when ensuring index template: %w", err)
	}
	return nil
}

// Index takes an Event an appends it to the stream
func (api API) Index(ctx context.Context, stream Stream, event Event) error {
	err := opensearch.IndexEvent(ctx, api.client, stream, event)
	if err != nil {
		return NewAPIError(err, "unable to index event")
	}
	return nil
}
