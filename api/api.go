package api

import (
	"context"
	"fmt"

	"github.com/kstiehl/index-bouncer/pkg/opensearch"
)

type (
	Stream = opensearch.DataStream
	Client = opensearch.Client
)

type API struct {
	client Client
}

// NewAPI returns an API object with the given configuration
func NewAPI(stream Stream, client Client) API {
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
