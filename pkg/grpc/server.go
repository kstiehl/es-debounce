package grpc

import (
	"context"
	"net"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/julienschmidt/httprouter"
	"github.com/kstiehl/index-bouncer/pkg/grpc/pb"
	"google.golang.org/grpc"
)

// And Option which can be applied to Options.
type Option = func(option *Options)

// WithListenAddress the default listening address.
func WithListenAddress(listen string) Option {
	return func(options *Options) {
		options.ListenAddress = listen
	}
}

// WithListen allow to directly configure a net.Listen for the server.
func WithListen(listener net.Listener) Option {
	return func(options *Options) {
		options.Listen = listener
	}
}

type Options struct {
	// Listen can be given to directly configure the port the grpc server is listening on.
	Listen net.Listener

	// ListenAddress can be used to configure a ListenAddress which is for the grpc server.
	// This will be ignored when Options.Listen is set.
	ListenAddress string
}

// InitDefaults initialises Options with default values for each setting.
func (o *Options) InitWithDefaults() {
	o.ListenAddress = ":8080"
	o.Listen = nil
}

// ApplyOptions iterates over []Option and applies every single one of them.
func (o *Options) ApplyOptions(options []Option) {
	for _, op := range options {
		op(o)
	}
}

// RunServer runs the server and block the goroutine.
func RunServer(ctx context.Context, options ...Option) error {
	log := logr.FromContextOrDiscard(ctx)

	serverOptions := Options{}
	serverOptions.InitWithDefaults()
	serverOptions.ApplyOptions(options)

	gServer := grpc.NewServer()

	recordingService, err := newRecordService(log.WithName("RecordingService"))
	if err != nil {
		log.Error(err, "unable to create recording service")
		return err
	}
	pb.RegisterRecordingServiceServer(gServer, &recordingService)

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
