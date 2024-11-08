package cmd

import (
	"github.com/alexkazantsev/experiments/rest/internal/user"
	"github.com/alexkazantsev/experiments/rest/server"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func RunServer() *cobra.Command {
	return &cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, _ []string) {
			fx.New(
				fx.Provide(
					server.NewRouter,
					server.NewServer,
				),
				user.Module,
				fx.Invoke(
					server.Run,
				),
			)
		},
	}
}
