package db

import (
	"database/sql"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type DB struct {
	db *sql.DB
}

func New() *DB {
	db, err := sql.Open("sqlite", filepath.Join(".", "app.db"))
	if err != nil {
		log.Fatal(err)
	}

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatal(err)

	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"sqlite", driver)
	if err != nil {
		log.Fatal(err)
	}
	if err = m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

	return &DB{db}
}

func (d *DB) Close() {
	d.db.Close()
}
