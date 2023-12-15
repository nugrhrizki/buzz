package whatsapp_user

import (
	"database/sql"
	"errors"

	"go.mau.fi/whatsmeow/types"

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

func (r *Repository) CreateUser(user *WhatsappUser) error {
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
		"INSERT INTO whatsapp_users (name, token) VALUES (?, ?)",
		user.Name,
		user.Token,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetConnectedUser() ([]WhatsappUser, error) {
	var users []WhatsappUser
	err := r.db.Select(
		&users,
		"SELECT * FROM whatsapp_users WHERE connected = 1",
	)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetUsers() ([]WhatsappUser, error) {
	var users []WhatsappUser
	err := r.db.Select(
		&users,
		"SELECT * FROM whatsapp_users",
	)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) GetUserById(id int) (*WhatsappUser, error) {
	var user WhatsappUser
	err := r.db.Get(
		&user,
		"SELECT * FROM whatsapp_users WHERE id = ?",
		id,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByJid(jid string) (*WhatsappUser, error) {
	var user WhatsappUser
	err := r.db.Get(
		&user,
		"SELECT * FROM whatsapp_users WHERE jid = ?",
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByToken(token string) (*WhatsappUser, error) {
	var user WhatsappUser
	err := r.db.Get(
		&user,
		"SELECT * FROM whatsapp_users WHERE token = ?",
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
		"SELECT qrcode FROM whatsapp_users WHERE id = ?",
		id,
	)
	if err != nil {
		return "", err
	}
	return qrcode, nil
}

func (r *Repository) SetUserConnected(id int, connected int) error {
	_, err := r.db.Exec(
		"UPDATE whatsapp_users SET connected = ? WHERE id = ?",
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
		"UPDATE whatsapp_users SET jid = ? WHERE id = ?",
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
		"UPDATE whatsapp_users SET qrcode = ? WHERE id = ?",
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
		"UPDATE whatsapp_users SET events = ? WHERE id = ?",
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
		"UPDATE whatsapp_users SET webhook = ? WHERE id = ?",
		webhook,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
