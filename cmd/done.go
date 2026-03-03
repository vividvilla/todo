package cmd

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vivek/todod/db"
)

var doneCmd = &cobra.Command{
	Use:   "done <id>",
	Short: "Mark a TODO as done",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid ID: %s", args[0])
		}

		if err := db.MarkDone(id); err != nil {
			return err
		}

		color.Green("✓ TODO #%d marked as done", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
