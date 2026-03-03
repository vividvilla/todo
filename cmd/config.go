package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vivek/todod/db"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := db.SetConfig(args[0], args[1]); err != nil {
			return fmt.Errorf("failed to set config: %w", err)
		}
		color.Green("✓ Config %s = %s", args[0], args[1])
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a config value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		val, err := db.GetConfig(args[0])
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}
		if val == "" {
			fmt.Printf("%s: (not set)\n", args[0])
		} else {
			fmt.Printf("%s: %s\n", args[0], val)
		}
		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	rootCmd.AddCommand(configCmd)
}
