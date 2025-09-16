package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  "Print the version number of the application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ec2-builder v1.0.0")
	},
}
