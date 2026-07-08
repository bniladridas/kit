package commands

import (
	"github.com/spf13/cobra"
)

func newVersionCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of kit",
		Long:  `Print the version number, build information, and platform details of kit.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("kit version " + version)
		},
	}
}
