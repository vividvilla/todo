package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vivek/todod/db"
	"github.com/vivek/todod/model"
)

var addCmd = &cobra.Command{
	Use:   "add <title>",
	Short: "Add a new TODO",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := strings.Join(args, " ")

		desc, _ := cmd.Flags().GetString("description")
		pri, _ := cmd.Flags().GetString("priority")
		due, _ := cmd.Flags().GetString("due")

		if !model.ValidPriority(pri) {
			return fmt.Errorf("invalid priority %q (use: low, medium, high)", pri)
		}

		var dueAt *time.Time
		if due != "" {
			parsed, err := parseDueDate(due)
			if err != nil {
				return fmt.Errorf("invalid due date %q: %w (use RFC3339 format, e.g. 2026-03-04T15:00:00Z)", due, err)
			}
			dueAt = &parsed
		}

		todo, err := db.AddTodo(title, desc, model.Priority(pri), model.StatusPending, dueAt)
		if err != nil {
			return fmt.Errorf("failed to add todo: %w", err)
		}

		color.Green("✓ Added TODO #%d: %s", todo.ID, todo.Title)
		if dueAt != nil {
			fmt.Printf("  Due: %s\n", dueAt.Local().Format("Mon Jan 2, 2006 3:04 PM"))
		}
		return nil
	},
}

func parseDueDate(s string) (time.Time, error) {
	// Try RFC3339
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	// Try date only
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	// Try datetime without timezone
	if t, err := time.ParseInLocation("2006-01-02T15:04:05", s, time.Local); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04", s, time.Local); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("unrecognized format")
}

func init() {
	addCmd.Flags().StringP("description", "d", "", "Description of the TODO")
	addCmd.Flags().StringP("priority", "p", "medium", "Priority: low, medium, high")
	addCmd.Flags().StringP("due", "D", "", "Due date/time (e.g. 2026-03-04T15:00:00Z or 2026-03-04)")
	rootCmd.AddCommand(addCmd)
}
