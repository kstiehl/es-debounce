package grpc

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/kstiehl/index-bouncer/pkg/grpc/pb"
)

type recordService struct {
	pb.UnimplementedRecordingServiceServer
	log logr.Logger
}

func newRecordService(log logr.Logger) (recordService, error) {
	return recordService{
		log: log,
	}, nil
}

func (r recordService) RecordEvent(_ context.Context, req *pb.RecordEventRequest) (*pb.RecordEventResponse, error) {
	r.log.V(3).Info("record event", "eventID", req.EventID, "objectID", req.ObjectID)
	return &pb.RecordEventResponse{Code: pb.RecordResponseCode_RECORD_OK}, nil
}
