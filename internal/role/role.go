package role

import "time"

type Role struct {
	Id        int64     `json:"id"         db:"id"`
	Name      string    `json:"name"       db:"name"`
	Actions   string    `json:"actions"    db:"actions"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" db:"deleted_at"`
}

func New() string {
	return `CREATE TABLE IF NOT EXISTS roles (
		id BIGSERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		actions TEXT NOT NULL DEFAULT '{}',
		created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ,
		deleted_at TIMESTAMPTZ
	);`
}
