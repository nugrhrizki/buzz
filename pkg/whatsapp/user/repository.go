package user

import (
	"database/sql"
	"errors"

	"go.mau.fi/whatsmeow/types"

	"github.com/nugrhrizki/buzz/pkg/database"
	"github.com/rs/zerolog"
)

type Repository struct {
	db  *database.Database
	log *zerolog.Logger
}

func NewRepository(db *database.Database, log *zerolog.Logger) *Repository {
	return &Repository{db, log}
}

func (r *Repository) Migration() string {
	return New()
}

func (r *Repository) CreateUser(user *User) error {
	_, err := r.GetUserByToken(user.Token)

	switch err {
	case sql.ErrNoRows:
		break
	case nil:
		return errors.New("user already exists")
	default:
		return err
	}

	_, err = r.db.Exec(
		"INSERT INTO whatsapp_users (name, token) VALUES ($1, $2)",
		user.Name,
		user.Token,
	)
	if err != nil {
		r.log.Error().Err(err).Msg("failed to create user")
		return err
	}
	return nil
}

func (r *Repository) UpdateUser(user *User) error {
	_, err := r.db.Exec(
		"UPDATE whatsapp_users SET name = $1, token = $2 WHERE id = $3",
		user.Name,
		user.Token,
		user.Id,
	)
	if err != nil {
		r.log.Error().Err(err).Msg("failed to update user")
		return err
	}
	return nil
}

func (r *Repository) DeleteUser(user *User) error {
	_, err := r.db.Exec(
		"DELETE FROM whatsapp_users WHERE id = $1",
		user.Id,
	)
	if err != nil {
		r.log.Error().Err(err).Msg("failed to delete user")
		return err
	}
	return nil
}

func (r *Repository) GetConnectedUser() ([]User, error) {
	var users []User
	err := r.db.Select(
		&users,
		"SELECT * FROM whatsapp_users WHERE connected = 1",
	)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetUsers() ([]User, error) {
	users := []User{}
	err := r.db.DB.Select(
		&users,
		"SELECT * FROM whatsapp_users ORDER BY id DESC",
	)
	if err != nil {
		r.log.Error().Err(err).Msg("failed to get users")
		return nil, err
	}

	return users, nil
}

func (r *Repository) GetUserById(id int) (*User, error) {
	var user User
	err := r.db.Get(
		&user,
		"SELECT * FROM whatsapp_users WHERE id = $1",
		id,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByJid(jid string) (*User, error) {
	var user User
	err := r.db.Get(
		&user,
		"SELECT * FROM whatsapp_users WHERE jid = $1",
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByToken(token string) (*User, error) {
	var user User
	err := r.db.Get(
		&user,
		"SELECT * FROM whatsapp_users WHERE token = $1",
		token,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetQRCode(id int) (string, error) {
	var qrcode string
	err := r.db.Get(
		&qrcode,
		"SELECT qrcode FROM whatsapp_users WHERE id = $1",
		id,
	)
	if err != nil {
		return "", err
	}
	return qrcode, nil
}

func (r *Repository) SetUserConnected(id int, connected int) error {
	_, err := r.db.Exec(
		"UPDATE whatsapp_users SET connected = $1 WHERE id = $2",
		connected,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) SetUserJid(id int, jid types.JID) error {
	_, err := r.db.Exec(
		"UPDATE whatsapp_users SET jid = $1 WHERE id = $2",
		jid,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) SetQRCode(id int, qrcode string) error {
	_, err := r.db.Exec(
		"UPDATE whatsapp_users SET qrcode = $1 WHERE id = $2",
		qrcode,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) SetEvents(id int, events string) error {
	_, err := r.db.Exec(
		"UPDATE whatsapp_users SET events = $1 WHERE id = $2",
		events,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) SetWebhook(id int, webhook string) error {
	_, err := r.db.Exec(
		"UPDATE whatsapp_users SET webhook = $1 WHERE id = $2",
		webhook,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
