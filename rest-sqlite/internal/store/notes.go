package store

import (
	"database/sql"
	"errors"
)

// ErrNotFound is returned when a row does not exist.
var ErrNotFound = errors.New("not found")

// Note is a row of the notes table.
type Note struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt string `json:"createdAt"`
}

// Notes is the repository for notes.
type Notes struct{ DB *sql.DB }

func (r Notes) Create(title, body string) (Note, error) {
	res, err := r.DB.Exec(`INSERT INTO notes (title, body) VALUES (?, ?)`, title, body)
	if err != nil {
		return Note{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Note{}, err
	}
	return r.Get(id)
}

func (r Notes) Get(id int64) (Note, error) {
	var n Note
	err := r.DB.QueryRow(
		`SELECT id, title, body, created_at FROM notes WHERE id = ?`, id,
	).Scan(&n.ID, &n.Title, &n.Body, &n.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Note{}, ErrNotFound
	}
	return n, err
}

func (r Notes) List(limit, offset int) ([]Note, error) {
	rows, err := r.DB.Query(
		`SELECT id, title, body, created_at FROM notes ORDER BY id DESC LIMIT ? OFFSET ?`,
		limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Note{}
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Body, &n.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

func (r Notes) Update(id int64, title, body string) (Note, error) {
	res, err := r.DB.Exec(`UPDATE notes SET title = ?, body = ? WHERE id = ?`, title, body, id)
	if err != nil {
		return Note{}, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return Note{}, ErrNotFound
	}
	return r.Get(id)
}

func (r Notes) Delete(id int64) error {
	res, err := r.DB.Exec(`DELETE FROM notes WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}
