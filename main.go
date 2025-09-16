package main

import (
	"log"

	cmd "github.com/bldmgr/ec2-builder/cmd"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	verbose bool
	config  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ec2-builder",
	Short: "Create new ec2 instances",
	Long:  `Create new ec2 instances`,
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "", "config file (default is $HOME/.cli.yaml)")

	// Add subcommands
	rootCmd.AddCommand(cmd.VersionCmd)
	rootCmd.AddCommand(cmd.CreateCmd)
	rootCmd.AddCommand(cmd.ListCmd)
	rootCmd.AddCommand(cmd.DeleteCmd)

	// Create command flags
	cmd.CreateCmd.Flags().StringP("image", "i", "", "image of resource to create (required)")
	cmd.CreateCmd.Flags().StringP("type", "t", "", "type of resource to create (required)")
	cmd.CreateCmd.Flags().StringP("name", "n", "", "name of the resource (required)")
	cmd.CreateCmd.MarkFlagRequired("type")
	cmd.CreateCmd.MarkFlagRequired("name")
	cmd.CreateCmd.MarkFlagRequired("image")

	// List command flags
	cmd.ListCmd.Flags().StringP("type", "t", "", "type of resource to list")

	// Delete command flags
	cmd.DeleteCmd.Flags().StringP("name", "n", "", "name of the resource to delete (required)")
	cmd.DeleteCmd.Flags().BoolP("force", "f", false, "force deletion without confirmation")
	cmd.DeleteCmd.MarkFlagRequired("name")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
