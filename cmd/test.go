package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func testCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Test service",
		Run: func(cmd *cobra.Command, args []string) {
			cmdd := exec.Command("go", "test", "./internal/test", "-v")
			cmdd.Stdout = os.Stdout
			cmdd.Stderr = os.Stderr

			if err := cmdd.Run(); err != nil {
				fmt.Println("Error running tests:", err)
				os.Exit(1)
			}
		},
	}
}
