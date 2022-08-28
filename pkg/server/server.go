package server

import (
	"context"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/julienschmidt/httprouter"
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

	log.WithValues("listenAddress", serverOptions.ListenAddress).Info("starting server")
	http.ListenAndServe(serverOptions.ListenAddress, assembleRouter())

	return nil
}

func assembleRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/_/ping", func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		w.Write([]byte("pong"))
	})
	return router
}
