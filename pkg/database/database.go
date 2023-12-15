package database

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nugrhrizki/buzz/pkg/env"
	"github.com/rs/zerolog"
)

type Database struct {
	DB  *sqlx.DB
	log *zerolog.Logger
}

type Seed struct {
	Schema string
	Data   interface{}
}

func NewDatabase(env *env.Env, log *zerolog.Logger) *Database {
	db, err := sqlx.Connect(
		env.GetDBDriver(),
		env.GetDBDSN(),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	return &Database{
		DB:  db,
		log: log,
	}
}

func (d *Database) Get(dest interface{}, query string, args ...interface{}) error {
	return d.DB.Get(dest, query, args...)
}

func (d *Database) Select(dest interface{}, query string, args ...interface{}) error {
	return d.DB.Select(dest, query, args...)
}

func (d *Database) Query(query string, args ...interface{}) (*sqlx.Rows, error) {
	return d.DB.Queryx(query, args...)
}

func (d *Database) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return d.DB.NamedQuery(query, arg)
}

func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.DB.Exec(query, args...)
}

func (d *Database) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return d.DB.NamedExec(query, arg)
}

func (d *Database) Begin() (*sql.Tx, error) {
	return d.DB.Begin()
}

func (d *Database) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (d *Database) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

func (d *Database) Close() {
	d.DB.Close()
}

type Migrate interface {
	Migration() string
}

func (d *Database) Migrate(migration ...Migrate) {
	tx, err := d.Begin()
	if err != nil {
		d.log.Fatal().Err(err).Msg("failed to begin transaction")
	}

	schemas := make([]string, len(migration))
	for i, m := range migration {
		schemas[i] = m.Migration()
	}

	for _, schema := range schemas {
		if _, err := d.DB.Exec(schema); err != nil {
			d.log.Warn().Err(err).Msg("failed to migrate")
		}
	}

	if err := d.Commit(tx); err != nil {
		d.Rollback(tx)
		d.log.Fatal().Err(err).Msg("failed to commit transaction")
	}
}

func (d *Database) Seeder(seeds []Seed) {
	tx, err := d.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for _, seed := range seeds {
		if _, err := d.DB.NamedExec(seed.Schema, seed.Data); err != nil {
			d.Rollback(tx)
			log.Fatal(err)
		}
	}

	if err := d.Commit(tx); err != nil {
		d.Rollback(tx)
		log.Fatal(err)
	}
}
