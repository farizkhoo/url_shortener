package main

import (
	"database/sql"

	"github.com/google/uuid"
)

// URL holds the uuid and the name of the Long
type URL struct {
	ID      string
	LongURL string
}

// GetURLFromUUID retrieves a long url from a uuid
func GetURLFromUUID(db *sql.DB, uuid string) (*URL, error) {
	var url URL
	err := db.QueryRow("SELECT uuid, long_url FROM url WHERE uuid = $1", uuid).Scan(&url.ID, &url.LongURL)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

// GetURLFromLongURL retrieves url from a long url
func GetURLFromLongURL(db *sql.DB, longURL string) (*URL, error) {
	var url URL
	err := db.QueryRow("SELECT uuid, long_url FROM url WHERE long_url = $1", longURL).Scan(&url.ID, &url.LongURL)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

// StoreURL accepts a db object and url and stores it into the database
func StoreURL(db *sql.DB, url string) error {
	uuid := uuid.New()
	_, err := db.Exec(`
		INSERT INTO url (
			uuid,
			long_url
		)
		VALUES ($1, $2)
	`, uuid, url)
	if err != nil {
		return err
	}
	return nil
}
