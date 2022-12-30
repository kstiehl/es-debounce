package grpc

import (
	"context"

	"github.com/kstiehl/index-bouncer/pkg/grpc/pb"
)

type streamingService struct {
	pb.UnimplementedStreamingServiceServer
}

func (s streamingService) Index(_ context.Context, _ *pb.Event) (*pb.IndexResonse, error) {
	panic("not implemented") // TODO: Implement
}
