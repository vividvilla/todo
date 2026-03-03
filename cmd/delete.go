package cmd

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vivek/todod/db"
)

var deleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Short:   "Delete a TODO",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid ID: %s", args[0])
		}

		if err := db.DeleteTodo(id); err != nil {
			return err
		}

		color.Yellow("✗ Deleted TODO #%d", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
