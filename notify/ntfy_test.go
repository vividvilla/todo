package notify

import (
	"strings"
	"testing"
	"time"

	"github.com/vivek/todod/model"
)

func TestMapPriority(t *testing.T) {
	tests := []struct {
		in  model.Priority
		out string
	}{
		{model.PriorityHigh, "5"},
		{model.PriorityMedium, "3"},
		{model.PriorityLow, "2"},
	}
	for _, tc := range tests {
		if got := mapPriority(tc.in); got != tc.out {
			t.Errorf("mapPriority(%q) = %q, want %q", tc.in, got, tc.out)
		}
	}
}

func TestMapTags(t *testing.T) {
	if got := mapTags(model.PriorityHigh); !strings.Contains(got, "rotating_light") {
		t.Errorf("expected rotating_light tag for high, got %q", got)
	}
	if got := mapTags(model.PriorityLow); !strings.Contains(got, "information_source") {
		t.Errorf("expected information_source tag for low, got %q", got)
	}
	if got := mapTags(model.PriorityMedium); !strings.Contains(got, "todo") {
		t.Errorf("expected todo tag, got %q", got)
	}
}

func TestFormatNtfyMessage(t *testing.T) {
	due := time.Date(2026, 3, 5, 14, 0, 0, 0, time.UTC)
	todo := model.Todo{
		ID:          1,
		Title:       "Test task",
		Description: "Some details",
		Priority:    model.PriorityHigh,
		Status:      model.StatusPending,
		DueAt:       &due,
	}

	msg := formatNtfyMessage(todo)

	if !strings.Contains(msg, "**Test task**") {
		t.Error("expected bold title in message")
	}
	if !strings.Contains(msg, "Some details") {
		t.Error("expected description in message")
	}
	if !strings.Contains(msg, "high") {
		t.Error("expected priority in message")
	}
	if !strings.Contains(msg, "Due:") {
		t.Error("expected due date in message")
	}
}

func TestFormatNtfyMessageNoDue(t *testing.T) {
	todo := model.Todo{
		ID:       2,
		Title:    "No due task",
		Priority: model.PriorityLow,
		Status:   model.StatusPending,
	}

	msg := formatNtfyMessage(todo)

	if strings.Contains(msg, "Due:") {
		t.Error("expected no due date line when DueAt is nil")
	}
}
