package grpc

import (
	context "context"

	"github.com/go-logr/logr"
	"github.com/kstiehl/index-bouncer/api"
	"github.com/kstiehl/index-bouncer/grpc/types"
)

type Server struct {
	types.UnimplementedStreamingServiceServer
}

func (Server) Index(ctx context.Context, event *types.Event) (*types.IndexResonse, error) {
	log := logr.FromContextOrDiscard(ctx).V(1).WithName("Indexer")
	if event == nil {
		log.Info("empty event received. Check client implementation")
	}

	// create new logger context so that log messages from now on contain the event.
	log = log.WithValues("eventID", event.GetEventID(),
		"objectID", event.ObjectID)
	ctx = logr.NewContext(ctx, log)

	api.Index(ctx, nil, event)
	return &types.IndexResonse{Code: types.StatusCode_RECORD_OK}, nil
}

func (Server) mustEmbedUnimplementedStreamingServiceServer() {
	panic("not implemented")
}
