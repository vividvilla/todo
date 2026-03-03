package notify

import (
	"fmt"
	"time"
)

// doWithRetry attempts fn up to 2 times (1 retry after 2s).
func doWithRetry(fn func() error) error {
	var lastErr error
	for attempt := 1; attempt <= 2; attempt++ {
		if attempt > 1 {
			time.Sleep(2 * time.Second)
		}
		if err := fn(); err != nil {
			lastErr = fmt.Errorf("attempt %d: %w", attempt, err)
			continue
		}
		return nil
	}
	return lastErr
}
