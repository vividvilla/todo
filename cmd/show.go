package cmd

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vivek/todod/db"
)

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show details of a TODO",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid ID: %s", args[0])
		}

		todo, err := db.GetTodo(id)
		if err != nil {
			return fmt.Errorf("todo #%d not found", id)
		}

		bold := color.New(color.Bold)

		bold.Printf("TODO #%d\n", todo.ID)
		fmt.Println("─────────────────────────────────")
		fmt.Printf("  Title:       %s\n", todo.Title)
		fmt.Printf("  Description: %s\n", defaultStr(todo.Description, "-"))
		fmt.Printf("  Priority:    %s\n", colorPriority(todo.Priority))
		fmt.Printf("  Status:      %s\n", colorStatus(todo.Status))

		if todo.DueAt != nil {
			fmt.Printf("  Due:         %s\n", todo.DueAt.Local().Format("Mon Jan 2, 2006 3:04 PM"))
		} else {
			fmt.Printf("  Due:         -\n")
		}

		fmt.Printf("  Notified:    %v\n", todo.Notified)
		fmt.Printf("  Created:     %s\n", todo.CreatedAt.Local().Format("Mon Jan 2, 2006 3:04 PM"))
		fmt.Printf("  Updated:     %s\n", todo.UpdatedAt.Local().Format("Mon Jan 2, 2006 3:04 PM"))
		return nil
	},
}

func defaultStr(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func init() {
	rootCmd.AddCommand(showCmd)
}
