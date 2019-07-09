package app

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // import postgres driver
	"github.com/pkg/errors"
)

// Postgres holds postgres login information
type Postgres struct {
	DBName   string
	User     string
	Host     string
	Password string
	DB       *sql.DB
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
func (p *Postgres) Connect() error {
	connStr := fmt.Sprintf(
		"user=%v dbname=%v host=%v password=%v sslmode=disable",
		p.User,
		p.DBName,
		p.Host,
		p.Password,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return errors.Wrap(err, "Failed to open database")
	}
	p.DB = db
	return nil
}

// GetURLFromUUID retrieves a long url from a uuid
func (p Postgres) GetURLFromUUID(uuid string) (*URL, error) {
	var url URL
	err := p.DB.QueryRow("SELECT uuid, long_url FROM url WHERE uuid = $1", uuid).Scan(&url.ID, &url.LongURL)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

// GetURLFromLongURL retrieves url from a long url
func (p Postgres) GetURLFromLongURL(longURL string) (*URL, error) {
	var url URL
	err := p.DB.QueryRow("SELECT uuid, long_url FROM url WHERE long_url = $1", longURL).Scan(&url.ID, &url.LongURL)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

// StoreURL accepts a db object and url and stores it into the database
func (p Postgres) StoreURL(url string) (string, error) {
	id := uuid.New()
	_, err := p.DB.Exec(`
		INSERT INTO url (
			uuid,
			long_url
		)
		VALUES ($1, $2)
	`, id, url)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
