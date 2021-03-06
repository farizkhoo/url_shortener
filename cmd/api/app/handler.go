package app

import (
	"database/sql"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

const (
	ip           string = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	urlschema    string = `((ftp|tcp|udp|wss?|https?):\/\/)`
	urlusername  string = `(\S+(:\S*)?@)`
	urlpath      string = `((\/|\?|#)[^\s]*)`
	urlport      string = `(:(\d{1,5}))`
	urlip        string = `([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`
	urlsubdomain string = `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`
	urlregex     string = `^` + urlschema + `?` + urlusername + `?` + `((` + urlip + `|(\[` + ip + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + urlsubdomain + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + urlport + `?` + urlpath + `?$`
)
const maxURLRuneCount = 2083
const minURLRuneCount = 3

var (
	rxURL = regexp.MustCompile(urlregex)
)

// IsURL checks if the string is a URL.
func IsURL(str string) bool {
	if str == "" || utf8.RuneCountInString(str) >= maxURLRuneCount || len(str) <= minURLRuneCount || strings.HasPrefix(str, ".") {
		return false
	}
	strTemp := str
	if strings.Contains(str, ":") && !strings.Contains(str, "://") {
		strTemp = "http://" + str
	}
	u, err := url.Parse(strTemp)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}
	return rxURL.MatchString(str)
}

type urlRepository interface {
	GetURLFromUUID(uuid string) (*URL, error)
	GetURLFromLongURL(longURL string) (*URL, error)
	StoreURL(longURL string) (string, error)
}

type getURLRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type createURLRequest struct {
	LongURL string `json:"url" binding:"required"`
}

// URLHandler has a urlRepository
type URLHandler struct {
	ur urlRepository
}

// NewURLHandler returns a new instance of url handler
func NewURLHandler(ur urlRepository) *URLHandler {
	return &URLHandler{ur}
}

// GetURL returns a gin context handling geturl requests
func (u *URLHandler) GetURL() func(c *gin.Context) {
	return func(c *gin.Context) {
		var req getURLRequest
		if err := c.ShouldBindUri(&req); err != nil {
			jsonCTX(c, http.StatusInternalServerError, "invalid url key provided", req.ID)
			return
		}

		url, err := u.ur.GetURLFromUUID(req.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				jsonCTX(c, http.StatusNotFound, "no result for url key provided", req.ID)
				return
			}
			jsonCTX(c, http.StatusInternalServerError, err.Error(), req.ID)
			return
		}
		c.Redirect(http.StatusMovedPermanently, url.LongURL)
	}
}

// CreateURL returns a gin context handling CreateURL requests
func (u *URLHandler) CreateURL() func(c *gin.Context) {
	return func(c *gin.Context) {
		var req createURLRequest

		if err := c.Bind(&req); err != nil {
			jsonCTX(c, http.StatusBadRequest, "failed to bind payload url", req.LongURL)
			return
		}

		if !IsURL(req.LongURL) {
			jsonCTX(c, http.StatusBadRequest, "invalid url string provided", req.LongURL)
			return
		}

		url, _ := u.ur.GetURLFromLongURL(req.LongURL)
		if url != nil {
			jsonCTX(c, http.StatusConflict, "url already exists in database", url.ID)
			return
		}

		uuid, err := u.ur.StoreURL(req.LongURL)
		if err != nil {
			jsonCTX(c, http.StatusInternalServerError, "failed to store url", req.LongURL)
			return
		}
		url, err = u.ur.GetURLFromUUID(uuid)
		if err != nil {
			jsonCTX(c, http.StatusInternalServerError, err.Error(), req.LongURL)
			return
		}

		jsonCTX(c, http.StatusOK, "successfully stored url", url.ID)
	}
}

func jsonCTX(c *gin.Context, status int, msg, value string) {
	c.JSON(status, gin.H{
		"msg":    msg,
		"status": status,
		"value":  value,
	})
}
