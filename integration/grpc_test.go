package integration

import (
	"context"
	"net"

	server "github.com/kstiehl/index-bouncer/pkg/grpc"
	"github.com/kstiehl/index-bouncer/pkg/grpc/pb"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
)

var _ = Describe("gRPC", Ordered, func() {
	var client pb.RecordingServiceClient
	BeforeAll(func() {
		listener, err := net.Listen("tcp", ":8080")
		Expect(err).ToNot(HaveOccurred())
		go server.RunServer(context.Background(), server.WithListen(listener))
		// TODO(kstiehl): find a way to stop server

		con, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
		client = pb.NewRecordingServiceClient(con)
	})

	It("Index document", func() {
		resp, err := client.RecordEvent(context.Background(), &pb.RecordEventRequest{
			EventID:  "testID",
			ObjectID: "objectID",
			Data:     "Data",
		})

		Expect(err).ToNot(HaveOccurred())
		Expect(resp).ToNot(BeNil())
		Expect(resp.Code).To(BeEquivalentTo(pb.RecordResponseCode_RECORD_OK))
	})

	It("Index multiple documents", func() {})
})
