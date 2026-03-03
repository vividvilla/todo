package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/vivek/todod/model"
)

func AddTodo(title, description string, priority model.Priority, status model.Status, dueAt *time.Time) (*model.Todo, error) {
	now := time.Now()
	var dueVal interface{}
	if dueAt != nil {
		dueVal = dueAt.UTC().Format(time.RFC3339)
	}

	res, err := DB.Exec(
		`INSERT INTO todos (title, description, priority, status, due_at, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		title, description, string(priority), string(status), dueVal,
		now.UTC().Format(time.RFC3339), now.UTC().Format(time.RFC3339),
	)
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()
	return &model.Todo{
		ID:          id,
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      status,
		DueAt:       dueAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func GetTodo(id int64) (*model.Todo, error) {
	row := DB.QueryRow(`SELECT id, title, description, priority, status, due_at, notified, created_at, updated_at FROM todos WHERE id = ?`, id)
	return scanTodo(row)
}

func ListTodos(status string, priority string, showAll bool) ([]model.Todo, error) {
	query := `SELECT id, title, description, priority, status, due_at, notified, created_at, updated_at FROM todos WHERE 1=1`
	var args []interface{}

	if !showAll && status == "" {
		query += ` AND status != 'done'`
	}
	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}
	if priority != "" {
		query += ` AND priority = ?`
		args = append(args, priority)
	}

	query += ` ORDER BY
		CASE priority WHEN 'high' THEN 1 WHEN 'medium' THEN 2 WHEN 'low' THEN 3 END,
		CASE WHEN due_at IS NOT NULL THEN 0 ELSE 1 END,
		due_at ASC,
		created_at DESC`

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		t, err := scanTodoRows(rows)
		if err != nil {
			return nil, err
		}
		todos = append(todos, *t)
	}
	return todos, nil
}

func UpdateTodo(id int64, title, description *string, priority *string, status *string, dueAt *string) error {
	// Build dynamic update
	sets := []string{}
	args := []interface{}{}

	if title != nil {
		sets = append(sets, "title = ?")
		args = append(args, *title)
	}
	if description != nil {
		sets = append(sets, "description = ?")
		args = append(args, *description)
	}
	if priority != nil {
		sets = append(sets, "priority = ?")
		args = append(args, *priority)
	}
	if status != nil {
		sets = append(sets, "status = ?")
		args = append(args, *status)
	}
	if dueAt != nil {
		if *dueAt == "" {
			sets = append(sets, "due_at = NULL")
		} else {
			sets = append(sets, "due_at = ?")
			args = append(args, *dueAt)
		}
	}

	if len(sets) == 0 {
		return fmt.Errorf("nothing to update")
	}

	sets = append(sets, "updated_at = ?")
	args = append(args, time.Now().UTC().Format(time.RFC3339))
	args = append(args, id)

	query := "UPDATE todos SET "
	for i, s := range sets {
		if i > 0 {
			query += ", "
		}
		query += s
	}
	query += " WHERE id = ?"

	res, err := DB.Exec(query, args...)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("todo #%d not found", id)
	}
	return nil
}

func MarkDone(id int64) error {
	s := "done"
	return UpdateTodo(id, nil, nil, nil, &s, nil)
}

func DeleteTodo(id int64) error {
	res, err := DB.Exec(`DELETE FROM todos WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("todo #%d not found", id)
	}
	return nil
}

func GetOverdueTodos() ([]model.Todo, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	rows, err := DB.Query(
		`SELECT id, title, description, priority, status, due_at, notified, created_at, updated_at
		 FROM todos
		 WHERE due_at IS NOT NULL AND due_at <= ? AND notified = 0 AND status != 'done'
		 ORDER BY due_at ASC`, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		t, err := scanTodoRows(rows)
		if err != nil {
			return nil, err
		}
		todos = append(todos, *t)
	}
	return todos, nil
}

func MarkNotified(id int64) error {
	_, err := DB.Exec(`UPDATE todos SET notified = 1, updated_at = ? WHERE id = ?`,
		time.Now().UTC().Format(time.RFC3339), id)
	return err
}

// scanners

func scanTodo(row *sql.Row) (*model.Todo, error) {
	var t model.Todo
	var dueAt, createdAt, updatedAt sql.NullString
	var priority, status string

	err := row.Scan(&t.ID, &t.Title, &t.Description, &priority, &status, &dueAt, &t.Notified, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	t.Priority = model.Priority(priority)
	t.Status = model.Status(status)

	if dueAt.Valid {
		parsed, _ := time.Parse(time.RFC3339, dueAt.String)
		t.DueAt = &parsed
	}
	if createdAt.Valid {
		t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt.String)
	}
	if updatedAt.Valid {
		t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt.String)
	}

	return &t, nil
}

func scanTodoRows(rows *sql.Rows) (*model.Todo, error) {
	var t model.Todo
	var dueAt, createdAt, updatedAt sql.NullString
	var priority, status string

	err := rows.Scan(&t.ID, &t.Title, &t.Description, &priority, &status, &dueAt, &t.Notified, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	t.Priority = model.Priority(priority)
	t.Status = model.Status(status)

	if dueAt.Valid {
		parsed, _ := time.Parse(time.RFC3339, dueAt.String)
		t.DueAt = &parsed
	}
	if createdAt.Valid {
		t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt.String)
	}
	if updatedAt.Valid {
		t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt.String)
	}

	return &t, nil
}
