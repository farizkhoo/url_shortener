package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	validUUID = "5537346a-84d0-472c-a466-8fd9af9c7db2"
	validURL  = "https://testurl.com"
)

type urlResponse struct {
	Message string `json:"msg"`
	Status  int    `json:"status"`
	Value   string `json:"value"`
}

type validURLRepository struct{}

func (t validURLRepository) GetURLFromUUID(uuid string) (*URL, error) {
	return &URL{
		ID:      uuid,
		LongURL: validURL,
	}, nil
}
func (t validURLRepository) GetURLFromLongURL(longURL string) (*URL, error) {
	return &URL{
		ID:      validUUID,
		LongURL: longURL,
	}, nil
}
func (t validURLRepository) StoreURL(longURL string) (string, error) {
	return validUUID, nil
}

func TestGetURLRoute(t *testing.T) {
	// GIVEN
	var vr validURLRepository
	uh := NewURLHandler(vr)
	router := SetupRouter(uh)
	w := httptest.NewRecorder()
	key := validUUID
	url := fmt.Sprintf("shorten_url/%s", key)
	req, _ := http.NewRequest("GET", url, nil)

	// WHEN
	router.ServeHTTP(w, req)

	// THEN
	assert.Equal(t, http.StatusMovedPermanently, w.Code)
}

type noDataURLRepository struct{}

func (t noDataURLRepository) GetURLFromUUID(uuid string) (*URL, error) {
	return &URL{
		ID:      uuid,
		LongURL: validURL,
	}, nil
}
func (t noDataURLRepository) GetURLFromLongURL(longURL string) (*URL, error) {
	return nil, sql.ErrNoRows
}
func (t noDataURLRepository) StoreURL(longURL string) (string, error) {
	return validUUID, nil
}

func TestCreateURLRoute(t *testing.T) {
	// GIVEN
	var ndr noDataURLRepository
	uh := NewURLHandler(ndr)
	router := SetupRouter(uh)
	w := httptest.NewRecorder()
	longURL := validURL
	payload := fmt.Sprintf(`{"url": "%s"}`, longURL)
	req, _ := http.NewRequest("POST", "shorten_url/", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	// WHEN
	router.ServeHTTP(w, req)

	// THEN
	var uresp urlResponse
	raw, _ := ioutil.ReadAll(w.Body)
	json.Unmarshal(raw, &uresp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, validUUID, uresp.Value)
	assert.Equal(t, "successfully stored url", uresp.Message)
}
