package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vivek/todod/db"
	"github.com/vivek/todod/model"
)

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit a TODO",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid ID: %s", args[0])
		}

		var title, description, priority, status, due *string

		if cmd.Flags().Changed("title") {
			v, _ := cmd.Flags().GetString("title")
			title = &v
		}
		if cmd.Flags().Changed("description") {
			v, _ := cmd.Flags().GetString("description")
			description = &v
		}
		if cmd.Flags().Changed("priority") {
			v, _ := cmd.Flags().GetString("priority")
			if !model.ValidPriority(v) {
				return fmt.Errorf("invalid priority %q", v)
			}
			priority = &v
		}
		if cmd.Flags().Changed("status") {
			v, _ := cmd.Flags().GetString("status")
			if !model.ValidStatus(v) {
				return fmt.Errorf("invalid status %q", v)
			}
			status = &v
		}
		if cmd.Flags().Changed("due") {
			v, _ := cmd.Flags().GetString("due")
			if v == "" {
				due = &v // clear due date
			} else {
				parsed, err := parseDueDate(v)
				if err != nil {
					return fmt.Errorf("invalid due date %q: %w", v, err)
				}
				formatted := parsed.UTC().Format(time.RFC3339)
				due = &formatted
			}
		}

		if err := db.UpdateTodo(id, title, description, priority, status, due); err != nil {
			return err
		}

		color.Green("✓ Updated TODO #%d", id)
		return nil
	},
}

func init() {
	editCmd.Flags().StringP("title", "t", "", "New title")
	editCmd.Flags().StringP("description", "d", "", "New description")
	editCmd.Flags().StringP("priority", "p", "", "New priority")
	editCmd.Flags().StringP("status", "s", "", "New status")
	editCmd.Flags().StringP("due", "D", "", "New due date (empty to clear)")
	rootCmd.AddCommand(editCmd)
}
