package main

import (
	"flag"
	"log"

	"github.com/alexkazantsev/experiments/rest/cmd"
	"github.com/spf13/pflag"

	"github.com/spf13/cobra"
)

func main() {
	var root = cobra.Command{Use: "root"}

	root.AddCommand(cmd.RunServer())

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	if err := root.Execute(); err != nil {
		log.Fatalf("fatal gateway execute error: %v", err)
	}
}
