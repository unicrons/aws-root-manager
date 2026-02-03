package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev"

func Version() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Version:", version)
			return nil
		},
	}
}
