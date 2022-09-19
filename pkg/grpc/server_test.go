package grpc

import (
	"context"
	"net"
	"testing"

	"github.com/kstiehl/index-bouncer/pkg/grpc/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestOptions(t *testing.T) {
	t.Parallel()

	t.Run("ApplyOptions", func(t *testing.T) {
		// given
		options := Options{}
		assert.Empty(t, options.ListenAddress, "ListenAddress should be empty")

		customOption := func(o *Options) {
			o.ListenAddress = "test"
		}

		// then
		options.ApplyOptions([]Option{customOption})

		// verify
		assert.Equal(t, "test", options.ListenAddress)
	})

	t.Run("DefaultInit", func(t *testing.T) {
		// given
		options := Options{}
		assert.Empty(t, options.ListenAddress, "ListenAddress should be empty")

		// then
		options.InitWithDefaults()

		// verify
		assert.Equal(t, ":8080", options.ListenAddress)
	})
}

func TestServer(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	listener, err := net.Listen("tcp", ":8080")
	assert.NoError(t, err)

	go RunServer(ctx, WithListen(listener))
	t.Run("RecordEvent", func(t *testing.T) {
		con, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
		assert.NoError(t, err)
		defer con.Close()
		client := pb.NewRecordingServiceClient(con)

		resp, err := client.RecordEvent(ctx, &pb.RecordEventRequest{EventID: "testID", ObjectID: "ObjectID", Data: "test"})
		assert.NoError(t, err)

		assert.Equal(t, pb.RecordResponseCode_RECORD_OK, resp.Code)
	})
}
