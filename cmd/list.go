package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/vivek/todod/db"
	"github.com/vivek/todod/model"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List TODOs",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		priority, _ := cmd.Flags().GetString("priority")
		showAll, _ := cmd.Flags().GetBool("all")

		if status != "" && !model.ValidStatus(status) {
			return fmt.Errorf("invalid status %q (use: pending, in-progress, done)", status)
		}
		if priority != "" && !model.ValidPriority(priority) {
			return fmt.Errorf("invalid priority %q (use: low, medium, high)", priority)
		}

		todos, err := db.ListTodos(status, priority, showAll)
		if err != nil {
			return fmt.Errorf("failed to list todos: %w", err)
		}

		if len(todos) == 0 {
			fmt.Println("No TODOs found.")
			return nil
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ID", "Title", "Priority", "Status", "Due"})

		for _, todo := range todos {
			dueStr := "-"
			if todo.DueAt != nil {
				dueStr = todo.DueAt.Local().Format("Jan 2, 2006 3:04 PM")
			}

			priStr := colorPriority(todo.Priority)
			statusStr := colorStatus(todo.Status)

			t.AppendRow(table.Row{
				todo.ID,
				truncate(todo.Title, 40),
				priStr,
				statusStr,
				dueStr,
			})
		}

		t.SetStyle(table.StyleRounded)
		t.Render()
		fmt.Printf("\nTotal: %d\n", len(todos))
		return nil
	},
}

func colorPriority(p model.Priority) string {
	switch p {
	case model.PriorityHigh:
		return color.RedString("high")
	case model.PriorityMedium:
		return color.YellowString("medium")
	case model.PriorityLow:
		return color.GreenString("low")
	}
	return string(p)
}

func colorStatus(s model.Status) string {
	switch s {
	case model.StatusDone:
		return color.GreenString("done")
	case model.StatusInProgress:
		return color.CyanString("in-progress")
	case model.StatusPending:
		return color.WhiteString("pending")
	}
	return string(s)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func init() {
	listCmd.Flags().StringP("status", "s", "", "Filter by status")
	listCmd.Flags().StringP("priority", "p", "", "Filter by priority")
	listCmd.Flags().BoolP("all", "a", false, "Show all including done")
	rootCmd.AddCommand(listCmd)
}
