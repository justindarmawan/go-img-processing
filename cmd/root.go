package cmd

import (
	"go-img-processing/bootstrap"

	"github.com/spf13/cobra"
)

func Execute() error {
	root := &cobra.Command{}
	config := bootstrap.Init()

	root.AddCommand(
		restCommand(config),
		testCommand(),
	)

	return root.Execute()
}
