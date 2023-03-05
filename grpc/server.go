package grpc

import (
	context "context"

	"github.com/go-logr/logr"
)

type Server struct {
	UnimplementedStreamingServiceServer
}

func (Server) Index(ctx context.Context, event *Event) (_ *IndexResonse, _ error) {
	log := logr.FromContextOrDiscard(ctx).V(1).WithName("Indexer")
	if event == nil {
		log.Info("empty event received. Check client implementation")
	}

	// create new logger context so that log messages from now on contain the event.
	log = log.WithValues("eventID", event.GetEventID(),
		"objectID", event.ObjectID)
	ctx = logr.NewContext(ctx, log)

	return &IndexResonse{Code: StatusCode_RECORD_OK}, nil
}

func (Server) mustEmbedUnimplementedStreamingServiceServer() {
	panic("not implemented")
}
