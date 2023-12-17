package role

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

func (r *Repository) CreateRole(role *Role) error {
	_, err := r.GetRoleByRolename(role.Name)

	switch err {
	case sql.ErrNoRows:
		break
	case nil:
		return errors.New("role already exists")
	default:
		return err
	}

	_, err = r.db.Exec(
		"INSERT INTO roles (name, actions) VALUES ($1, $2)",
		role.Name,
		role.Actions,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetRoleByRolename(rolename string) (*Role, error) {
	var role Role
	err := r.db.Get(
		&role,
		"SELECT * FROM roles WHERE name = $1 AND deleted_at IS NULL",
		rolename,
	)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repository) GetRoleById(id int) (*Role, error) {
	var role Role
	err := r.db.Get(
		&role,
		"SELECT * FROM roles WHERE id = $1 AND deleted_at IS NULL",
		id,
	)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repository) GetRoles() ([]Role, error) {
	var roles []Role
	err := r.db.Select(
		&roles,
		"SELECT * FROM roles WHERE deleted_at IS NULL",
	)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *Repository) UpdateRole(role *Role) error {
	_, err := r.GetRoleByRolename(role.Name)
	switch err {
	case sql.ErrNoRows:
		break
	case nil:
		return errors.New("role name already used")
	default:
		return err
	}

	_, err = r.db.Exec(
		`UPDATE roles
		SET
			name = $1,
			actions = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			id = $3`,
		role.Name,
		role.Actions,
		role.Id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteRole(role *Role) error {
	_, err := r.db.Exec(
		"UPDATE roles SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1",
		role.Id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) HardDeleteRole(role *Role) error {
	_, err := r.db.Exec(
		"DELETE FROM roles WHERE id = $1",
		role.Id,
	)
	if err != nil {
		return err
	}
	return nil
}
