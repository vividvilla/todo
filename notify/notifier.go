package notify

import "github.com/vivek/todod/model"

// Notifier is the interface every notification channel must implement.
type Notifier interface {
	// Name returns a human-readable channel name for logging.
	Name() string
	// Send delivers a notification for the given todo. Returns nil on success.
	Send(todo model.Todo) error
}
