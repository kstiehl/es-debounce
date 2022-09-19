package grpc

import (
	"context"
	"fmt"

	"github.com/kstiehl/index-bouncer/pkg/server/pb"
)

type recordService struct {
	pb.UnimplementedRecordingServiceServer
}

func newRecordService() (recordService, error) {
	return recordService{}, nil
}

func (recordService) RecordEvent(_ context.Context, _ *pb.RecordEventRequest) (*pb.RecordEventResponse, error) {
	fmt.Println("Log log log log")
	return &pb.RecordEventResponse{Code: pb.RecordResponseCode_RECORD_OK}, nil
}
