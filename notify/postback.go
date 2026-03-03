package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vivek/todod/model"
)

// PostbackNotifier sends a JSON POST to a generic webhook URL.
type PostbackNotifier struct {
	URL string
}

type PostbackPayload struct {
	Event     string     `json:"event"`
	Todo      model.Todo `json:"todo"`
	Timestamp string     `json:"timestamp"`
}

func NewPostback(url string) *PostbackNotifier {
	return &PostbackNotifier{URL: url}
}

func (p *PostbackNotifier) Name() string {
	return "postback"
}

func (p *PostbackNotifier) Send(todo model.Todo) error {
	payload := PostbackPayload{
		Event:     "todo_due",
		Todo:      todo,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	return doWithRetry(func() error {
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Post(p.URL, "application/json", bytes.NewReader(body))
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
