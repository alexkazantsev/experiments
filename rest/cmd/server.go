package cmd

import (
	"github.com/alexkazantsev/experiments/rest/server"
	"github.com/spf13/cobra"
)

func RunServer() *cobra.Command {
	return &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return server.Run()
		},
	}
}
