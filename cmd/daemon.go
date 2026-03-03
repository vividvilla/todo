package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vivek/todod/db"
	"github.com/vivek/todod/notify"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start the notification daemon",
	Long: `Starts a background daemon that checks for overdue TODOs every 60 seconds
and sends notifications via all configured channels.

Supported channels:
  postback  — Generic webhook (JSON POST)
  ntfy      — ntfy.sh push notifications

Configure one or more channels:
  todo config set postback_url https://example.com/webhook
  todo config set ntfy_url https://ntfy.sh/my_topic`,
	RunE: func(cmd *cobra.Command, args []string) error {
		interval, _ := cmd.Flags().GetInt("interval")

		// Write PID file
		pidFile := filepath.Join(db.DataDir(), "todo.pid")
		if err := os.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0644); err != nil {
			log.Printf("Warning: failed to write PID file: %v", err)
		}
		defer os.Remove(pidFile)

		color.Cyan("🔔 todo daemon started (checking every %ds)", interval)
		color.Cyan("   PID: %d", os.Getpid())
		color.Cyan("   PID file: %s", pidFile)
		color.Cyan("   DB: %s", db.DBPath())

		// Build notifiers from config
		notifiers := buildNotifiers()
		if len(notifiers) == 0 {
			color.Yellow("⚠  No notification channels configured. Set one or more:")
			color.Yellow("   todo config set postback_url <url>")
			color.Yellow("   todo config set ntfy_url <url>")
			color.Yellow("   Daemon will still run and check, but won't send notifications.")
		} else {
			for _, n := range notifiers {
				color.Cyan("   Channel: %s", n.Name())
			}
		}
		fmt.Println()

		// Set up signal handling
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		// Run immediately on start
		checkAndNotify()

		for {
			select {
			case <-ticker.C:
				checkAndNotify()
			case sig := <-sigCh:
				log.Printf("\nReceived %s, shutting down daemon...", sig)
				return nil
			}
		}
	},
}

// buildNotifiers reads config and returns all configured notifiers.
func buildNotifiers() []notify.Notifier {
	var notifiers []notify.Notifier

	if url, _ := db.GetConfig("postback_url"); url != "" {
		notifiers = append(notifiers, notify.NewPostback(url))
	}
	if url, _ := db.GetConfig("ntfy_url"); url != "" {
		notifiers = append(notifiers, notify.NewNtfy(url))
	}

	return notifiers
}

func checkAndNotify() {
	todos, err := db.GetOverdueTodos()
	if err != nil {
		log.Printf("Error fetching overdue todos: %v", err)
		return
	}

	if len(todos) == 0 {
		return
	}

	notifiers := buildNotifiers()

	for _, todo := range todos {
		log.Printf("📋 Overdue: #%d %q (due: %s)",
			todo.ID, todo.Title, todo.DueAt.Local().Format("Jan 2, 2006 3:04 PM"))

		if len(notifiers) == 0 {
			log.Printf("   ⚠ Skipping notification (no channels configured)")
			continue
		}

		// Fan-out to all channels. Mark notified if at least one succeeds.
		anySuccess := false
		for _, n := range notifiers {
			if err := n.Send(todo); err != nil {
				log.Printf("   ✗ [%s] Failed: %v", n.Name(), err)
			} else {
				log.Printf("   ✓ [%s] Sent successfully", n.Name())
				anySuccess = true
			}
		}

		if anySuccess {
			if err := db.MarkNotified(todo.ID); err != nil {
				log.Printf("   ✗ Failed to mark as notified: %v", err)
			}
		}
	}
}

func init() {
	daemonCmd.Flags().IntP("interval", "i", 60, "Check interval in seconds")
	rootCmd.AddCommand(daemonCmd)
}
