package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // import postgres driver
	"github.com/pkg/errors"
)

// Postgres holds postgres login information
type Postgres struct {
	DBName   string
	User     string
	Host     string
	Password string
}

// NewPostgres returns a new instance of postgres
func NewPostgres() *Postgres {
	return &Postgres{
		DBName:   os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Host:     os.Getenv("DB_HOST"),
		Password: os.Getenv("DB_PASSWORD"),
	}
}

// Connect returns a sql.DB object
func (p Postgres) Connect() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"user=%v dbname=%v host=%v password=%v sslmode=disable",
		p.User,
		p.DBName,
		p.Host,
		p.Password,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open database")
	}
	return db, nil
}
