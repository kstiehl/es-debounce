package cmd

import (
	"context"
	"log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"github.com/kstiehl/index-bouncer/pkg/server"
	"github.com/spf13/cobra"
)

func ServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "start the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := logr.NewContext(context.Background(),
				stdr.New(log.New(os.Stdout, "", log.LstdFlags)))
			return server.RunServer(ctx)
		},
	}

}
