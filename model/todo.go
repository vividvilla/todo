package model

import "time"

type Priority string
type Status string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
)

type Todo struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Priority    Priority   `json:"priority"`
	Status      Status     `json:"status"`
	DueAt       *time.Time `json:"due_at,omitempty"`
	Notified    bool       `json:"notified"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func ValidPriority(p string) bool {
	switch Priority(p) {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	}
	return false
}

func ValidStatus(s string) bool {
	switch Status(s) {
	case StatusPending, StatusInProgress, StatusDone:
		return true
	}
	return false
}
