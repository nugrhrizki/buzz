package env

type Env struct {
	DB_DRIVER   string
	DB_DSN      string
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	KeyLength   uint32
	SaltLength  uint32
	Secret      string
}

func New() *Env {
	return &Env{
		DB_DRIVER: "postgres",
		DB_DSN:    "user=postgres password=localdb dbname=db_whatsapp sslmode=disable",

		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		KeyLength:   32,
		SaltLength:  16,

		Secret: "secret",
	}
}
