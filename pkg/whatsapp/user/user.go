package user

type User struct {
	Id         int    `db:"id"         json:"id"`
	Name       string `db:"name"       json:"name"`
	Token      string `db:"token"      json:"token"`
	Webhook    string `db:"webhook"    json:"webhook"`
	Jid        string `db:"jid"        json:"jid"`
	Qrcode     string `db:"qrcode"     json:"qrcode"`
	Connected  *int   `db:"connected"  json:"connected"`
	Expiration *int   `db:"expiration" json:"expiration"`
	Events     string `db:"events"     json:"events"`
}

type UserInfo struct {
	Id      string `json:"id"`
	Jid     string `json:"jid"`
	Webhook string `json:"webhook"`
	Token   string `json:"token"`
	Events  string `json:"events"`
}

func New() string {
	return `CREATE TABLE IF NOT EXISTS whatsapp_users (
		id BIGSERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		token TEXT NOT NULL,
		webhook TEXT NOT NULL DEFAULT '',
		jid TEXT NOT NULL DEFAULT '',
		qrcode TEXT NOT NULL DEFAULT '',
		connected INTEGER,
		expiration INTEGER,
		events TEXT NOT NULL DEFAULT 'All'
	);`
}
