package notify

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vivek/todod/model"
)

// NtfyNotifier sends notifications to an ntfy topic URL.
// See https://docs.ntfy.sh/publish/
type NtfyNotifier struct {
	// TopicURL is the full ntfy topic URL, e.g. https://ntfy.sh/my_topic
	TopicURL string
}

func NewNtfy(topicURL string) *NtfyNotifier {
	return &NtfyNotifier{TopicURL: strings.TrimRight(topicURL, "/")}
}

func (n *NtfyNotifier) Name() string {
	return "ntfy"
}

func (n *NtfyNotifier) Send(todo model.Todo) error {
	title := fmt.Sprintf("⏰ TODO #%d Overdue", todo.ID)
	message := formatNtfyMessage(todo)
	priority := mapPriority(todo.Priority)
	tags := mapTags(todo.Priority)

	return doWithRetry(func() error {
		req, err := http.NewRequest("POST", n.TopicURL, strings.NewReader(message))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Title", title)
		req.Header.Set("Priority", priority)
		req.Header.Set("Tags", tags)
		req.Header.Set("Markdown", "yes")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		return fmt.Errorf("got status %d", resp.StatusCode)
	})
}

func formatNtfyMessage(todo model.Todo) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("**%s**\n", todo.Title))

	if todo.Description != "" {
		b.WriteString(fmt.Sprintf("%s\n", todo.Description))
	}

	b.WriteString(fmt.Sprintf("\n- **Priority:** %s\n", todo.Priority))
	b.WriteString(fmt.Sprintf("- **Status:** %s\n", todo.Status))

	if todo.DueAt != nil {
		b.WriteString(fmt.Sprintf("- **Due:** %s\n", todo.DueAt.Local().Format("Mon Jan 2, 2006 3:04 PM")))
	}

	return b.String()
}

// mapPriority converts todo priority to ntfy priority string.
// ntfy: 1=min, 2=low, 3=default, 4=high, 5=urgent
func mapPriority(p model.Priority) string {
	switch p {
	case model.PriorityHigh:
		return "5"
	case model.PriorityMedium:
		return "3"
	case model.PriorityLow:
		return "2"
	}
	return "3"
}

// mapTags returns ntfy emoji tags based on priority.
func mapTags(p model.Priority) string {
	switch p {
	case model.PriorityHigh:
		return "rotating_light,todo"
	case model.PriorityMedium:
		return "warning,todo"
	case model.PriorityLow:
		return "information_source,todo"
	}
	return "todo"
}
