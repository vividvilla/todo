package db

import (
	"os"
	"testing"
	"time"

	"github.com/vivek/todod/model"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	os.Setenv("HOME", dir)
	if err := Init(); err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	t.Cleanup(func() { Close() })
}

func TestAddAndGetTodo(t *testing.T) {
	setupTestDB(t)

	due := time.Now().Add(24 * time.Hour)
	todo, err := AddTodo("Buy milk", "From the store", model.PriorityHigh, model.StatusPending, &due)
	if err != nil {
		t.Fatalf("AddTodo failed: %v", err)
	}
	if todo.ID != 1 {
		t.Errorf("expected ID 1, got %d", todo.ID)
	}

	got, err := GetTodo(todo.ID)
	if err != nil {
		t.Fatalf("GetTodo failed: %v", err)
	}
	if got.Title != "Buy milk" {
		t.Errorf("expected title 'Buy milk', got %q", got.Title)
	}
	if got.Priority != model.PriorityHigh {
		t.Errorf("expected priority high, got %q", got.Priority)
	}
	if got.DueAt == nil {
		t.Error("expected due date to be set")
	}
}

func TestListTodos(t *testing.T) {
	setupTestDB(t)

	AddTodo("Task 1", "", model.PriorityHigh, model.StatusPending, nil)
	AddTodo("Task 2", "", model.PriorityLow, model.StatusPending, nil)
	AddTodo("Task 3", "", model.PriorityMedium, model.StatusDone, nil)

	// Default: hide done
	todos, err := ListTodos("", "", false)
	if err != nil {
		t.Fatalf("ListTodos failed: %v", err)
	}
	if len(todos) != 2 {
		t.Errorf("expected 2 active todos, got %d", len(todos))
	}

	// Show all
	todos, err = ListTodos("", "", true)
	if err != nil {
		t.Fatalf("ListTodos (all) failed: %v", err)
	}
	if len(todos) != 3 {
		t.Errorf("expected 3 total todos, got %d", len(todos))
	}

	// Filter by status
	todos, err = ListTodos("done", "", false)
	if err != nil {
		t.Fatalf("ListTodos (done) failed: %v", err)
	}
	if len(todos) != 1 {
		t.Errorf("expected 1 done todo, got %d", len(todos))
	}

	// Filter by priority
	todos, err = ListTodos("", "high", false)
	if err != nil {
		t.Fatalf("ListTodos (high) failed: %v", err)
	}
	if len(todos) != 1 {
		t.Errorf("expected 1 high-priority todo, got %d", len(todos))
	}
}

func TestUpdateTodo(t *testing.T) {
	setupTestDB(t)

	AddTodo("Original", "", model.PriorityLow, model.StatusPending, nil)

	newTitle := "Updated"
	newPri := "high"
	if err := UpdateTodo(1, &newTitle, nil, &newPri, nil, nil); err != nil {
		t.Fatalf("UpdateTodo failed: %v", err)
	}

	got, _ := GetTodo(1)
	if got.Title != "Updated" {
		t.Errorf("expected title 'Updated', got %q", got.Title)
	}
	if got.Priority != model.PriorityHigh {
		t.Errorf("expected priority high, got %q", got.Priority)
	}
}

func TestMarkDone(t *testing.T) {
	setupTestDB(t)

	AddTodo("Finish this", "", model.PriorityMedium, model.StatusPending, nil)
	if err := MarkDone(1); err != nil {
		t.Fatalf("MarkDone failed: %v", err)
	}

	got, _ := GetTodo(1)
	if got.Status != model.StatusDone {
		t.Errorf("expected status done, got %q", got.Status)
	}
}

func TestDeleteTodo(t *testing.T) {
	setupTestDB(t)

	AddTodo("Delete me", "", model.PriorityLow, model.StatusPending, nil)
	if err := DeleteTodo(1); err != nil {
		t.Fatalf("DeleteTodo failed: %v", err)
	}

	_, err := GetTodo(1)
	if err == nil {
		t.Error("expected error getting deleted todo")
	}
}

func TestGetOverdueTodos(t *testing.T) {
	setupTestDB(t)

	past := time.Now().Add(-1 * time.Hour)
	future := time.Now().Add(24 * time.Hour)

	AddTodo("Overdue", "", model.PriorityHigh, model.StatusPending, &past)
	AddTodo("Not yet", "", model.PriorityMedium, model.StatusPending, &future)
	AddTodo("No due", "", model.PriorityLow, model.StatusPending, nil)

	overdue, err := GetOverdueTodos()
	if err != nil {
		t.Fatalf("GetOverdueTodos failed: %v", err)
	}
	if len(overdue) != 1 {
		t.Errorf("expected 1 overdue todo, got %d", len(overdue))
	}
	if len(overdue) > 0 && overdue[0].Title != "Overdue" {
		t.Errorf("expected 'Overdue', got %q", overdue[0].Title)
	}
}

func TestMarkNotified(t *testing.T) {
	setupTestDB(t)

	past := time.Now().Add(-1 * time.Hour)
	AddTodo("Notify me", "", model.PriorityHigh, model.StatusPending, &past)

	if err := MarkNotified(1); err != nil {
		t.Fatalf("MarkNotified failed: %v", err)
	}

	// Should no longer appear in overdue
	overdue, _ := GetOverdueTodos()
	if len(overdue) != 0 {
		t.Errorf("expected 0 overdue after notified, got %d", len(overdue))
	}
}
