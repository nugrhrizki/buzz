package user

import "time"

type User struct {
	Id        int64      `json:"id"         db:"id"`
	Name      string     `json:"name"       db:"name"`
	Username  string     `json:"username"   db:"username"`
	Password  string     `json:"password"   db:"password"`
	Confirmed bool       `json:"confirmed"  db:"confirmed"`
	Whatsapp  *string    `json:"whatsapp"   db:"whatsapp"`
	Email     *string    `json:"email"      db:"email"`
	RoleId    int64      `json:"role_id"    db:"role_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

func New() string {
	return `CREATE TABLE IF NOT EXISTS users (
		id BIGSERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		confirmed BOOLEAN NOT NULL DEFAULT FALSE,
		whatsapp TEXT,
		email TEXT,
		role_id INTEGER NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ,
		deleted_at TIMESTAMPTZ
	);
	
	CREATE UNIQUE INDEX IF NOT EXISTS users_username_uindex ON users (username);
	ALTER TABLE users ADD CONSTRAINT users_role_id_fk FOREIGN KEY (role_id) REFERENCES roles (id);`
}
