package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

func testCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Test service",
		Run: func(cmd *cobra.Command, args []string) {
			cmdTest := exec.Command("sh", "-c", "go clean -testcache && go test ./internal/test -v")
			cmdTest.Stdout = cmd.OutOrStdout()
			cmdTest.Stderr = cmd.OutOrStderr()

			if err := cmdTest.Run(); err != nil {
				fmt.Println("error running tests:", err)
				return
			}
		},
	}
}
