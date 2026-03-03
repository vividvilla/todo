package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vivek/todod/db"
)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A CLI TODO tracker with notification daemon",
	Long: `todo is a command-line TODO tracker that stores tasks in SQLite
and includes a daemon for sending webhook notifications when tasks are due.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return db.Init()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		db.Close()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
