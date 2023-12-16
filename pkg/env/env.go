package env

type Env struct {
	DB_DRIVER string
	DB_DSN    string
}

func New() *Env {
	return &Env{
		DB_DRIVER: "postgres",
		DB_DSN:    "user=postgres password=postgres dbname=db_whatsapp sslmode=disable",
	}
}

func (e *Env) GetDBDriver() string {
	return e.DB_DRIVER
}

func (e *Env) GetDBDSN() string {
	return e.DB_DSN
}
