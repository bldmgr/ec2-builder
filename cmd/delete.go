package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a resource",
	Long:  "Delete a resource by name and type",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		force, _ := cmd.Flags().GetBool("force")

		if name == "" {
			fmt.Println("Error: name is required")
			cmd.Help()
			os.Exit(1)
		}

		if !force {
			fmt.Printf("Are you sure you want to delete '%s'? Use --force to skip confirmation.\n", name)
			return
		}

		fmt.Printf("Successfully deleted: %s\n", name)
	},
}
