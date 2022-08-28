package main

import (
	"fmt"
	"os"

	"github.com/kstiehl/debounce-indexer/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

func main() {
	rootCmd.AddCommand(cmd.ServeCmd())
	fmt.Println("starting")

	if err := rootCmd.Execute(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
