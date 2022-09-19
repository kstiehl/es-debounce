package grpc

import (
	"context"
	"net"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/julienschmidt/httprouter"
	"github.com/kstiehl/index-bouncer/pkg/server/pb"
	"google.golang.org/grpc"
)

// And Option which can be applied to Options.
type Option = func(option *Options)

type Options struct {
	ListenAddress string
}

// InitDefaults initialises Options with default values for each setting.
func (o *Options) InitWithDefaults() {
	o.ListenAddress = ":8080"
}

// ApplyOptions iterates over []Option and applies every single one of them.
func (o *Options) ApplyOptions(options []Option) {
	for _, op := range options {
		op(o)
	}
}

func RunServer(ctx context.Context, options ...Option) error {
	log := logr.FromContextOrDiscard(ctx)

	serverOptions := Options{}
	serverOptions.InitWithDefaults()
	serverOptions.ApplyOptions(options)

	gServer := grpc.NewServer()
	pb.RegisterRecordingServiceServer(gServer, nil)

	listen, err := net.Listen("tcp", serverOptions.ListenAddress)
	if err != nil {
		log.Error(err, "couldn't create listen", "port", serverOptions.ListenAddress)
	}
	log.Info("server listening on", "port", serverOptions.ListenAddress)

	if err := gServer.Serve(listen); err != nil {
		log.Error(err, "error when listening", "port", serverOptions.ListenAddress)
	}
	return nil
}

func assembleRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/_/ping", func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		w.Write([]byte("pong"))
	})
	return router
}
