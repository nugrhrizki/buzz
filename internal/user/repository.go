package user

import (
	"database/sql"
	"errors"

	"github.com/nugrhrizki/buzz/pkg/database"
)

type Repository struct {
	db *database.Database
}

func NewRepository(db *database.Database) *Repository {
	return &Repository{db}
}

func (r *Repository) Migration() string {
	return New()
}

func (r *Repository) CreateUser(user *User) error {
	_, err := r.GetUserByUsername(user.Username)

	switch err {
	case sql.ErrNoRows:
		break
	case nil:
		return errors.New("user already exists")
	default:
		return err
	}

	_, err = r.db.Exec(
		"INSERT INTO users (name, username, password, role_id) VALUES ($1, $2, $3, $4)",
		user.Name,
		user.Username,
		user.Password,
		user.RoleId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetUserByUsername(username string) (*User, error) {
	var user User
	err := r.db.Get(
		&user,
		"SELECT * FROM users WHERE username = $1 AND deleted_at IS NULL",
		username,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserById(id int) (*User, error) {
	var user User
	err := r.db.Get(
		&user,
		"SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL",
		id,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUsers() ([]User, error) {
	var users []User
	err := r.db.Select(
		&users,
		"SELECT * FROM users AND deleted_at IS NULL",
	)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) UpdateUser(user *User) error {
	_, err := r.GetUserByUsername(user.Username)
	switch err {
	case sql.ErrNoRows:
		break
	case nil:
		return errors.New("username already used")
	default:
		return err
	}

	_, err = r.db.Exec(
		`UPDATE users
		SET
			name = $1,
			username = $2,
			password = $3,
			confirmed = $4,
			whatsapp = $5,
			email = $6,
			role_id = $7
			updated_at = CURRENT_TIMESTAMP
		WHERE
			id = $8`,
		user.Name,
		user.Username,
		user.Password,
		user.Confirmed,
		user.Whatsapp,
		user.Email,
		user.RoleId,
		user.Id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteUser(user *User) error {
	_, err := r.db.Exec(
		"UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1",
		user.Id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) HardDeleteUser(user *User) error {
	_, err := r.db.Exec(
		"DELETE FROM users WHERE id = $1",
		user.Id,
	)
	if err != nil {
		return err
	}
	return nil
}
