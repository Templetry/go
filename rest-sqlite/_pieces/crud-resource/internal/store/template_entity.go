package store

import (
	"database/sql"
	"errors"
)

// init registers this resource's table through the migrations socket, so
// no existing file changes (ADR-0014).
func init() {
	Register(`
CREATE TABLE IF NOT EXISTS template_entities (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	name       TEXT NOT NULL,
	notes      TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL DEFAULT (datetime('now'))
);`)
}

// TemplateEntity is a row of the template_entities table.
type TemplateEntity struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Notes     string `json:"notes"`
	CreatedAt string `json:"createdAt"`
}

// TemplateEntities is the repository for template_entities.
type TemplateEntities struct{ DB *sql.DB }

func (r TemplateEntities) Create(name, notes string) (TemplateEntity, error) {
	res, err := r.DB.Exec(`INSERT INTO template_entities (name, notes) VALUES (?, ?)`, name, notes)
	if err != nil {
		return TemplateEntity{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return TemplateEntity{}, err
	}
	return r.Get(id)
}

func (r TemplateEntities) Get(id int64) (TemplateEntity, error) {
	var row TemplateEntity
	err := r.DB.QueryRow(
		`SELECT id, name, notes, created_at FROM template_entities WHERE id = ?`, id,
	).Scan(&row.ID, &row.Name, &row.Notes, &row.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return TemplateEntity{}, ErrNotFound
	}
	return row, err
}

func (r TemplateEntities) List(limit, offset int) ([]TemplateEntity, error) {
	rows, err := r.DB.Query(
		`SELECT id, name, notes, created_at FROM template_entities ORDER BY id DESC LIMIT ? OFFSET ?`,
		limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []TemplateEntity{}
	for rows.Next() {
		var row TemplateEntity
		if err := rows.Scan(&row.ID, &row.Name, &row.Notes, &row.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (r TemplateEntities) Update(id int64, name, notes string) (TemplateEntity, error) {
	res, err := r.DB.Exec(`UPDATE template_entities SET name = ?, notes = ? WHERE id = ?`, name, notes, id)
	if err != nil {
		return TemplateEntity{}, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return TemplateEntity{}, ErrNotFound
	}
	return r.Get(id)
}

func (r TemplateEntities) Delete(id int64) error {
	res, err := r.DB.Exec(`DELETE FROM template_entities WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}
