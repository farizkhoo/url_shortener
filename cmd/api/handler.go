package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// const (
// 	ip           string = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
// 	urlschema    string = `((ftp|tcp|udp|wss?|https?):\/\/)`
// 	urlusername  string = `(\S+(:\S*)?@)`
// 	urlpath      string = `((\/|\?|#)[^\s]*)`
// 	urlport      string = `(:(\d{1,5}))`
// 	urlip        string = `([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`
// 	urlsubdomain string = `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`
// 	urlregex     string = `^` + urlschema + `?` + urlusername + `?` + `((` + urlip + `|(\[` + ip + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + urlsubdomain + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + urlport + `?` + urlpath + `?$`
// )
// const maxURLRuneCount = 2083
// const minURLRuneCount = 3

// var (
// 	rxURL = regexp.MustCompile(urlregex)
// )

type request struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func jsonCTX(c *gin.Context, status int, msg, value string) {
	c.JSON(status, gin.H{
		"msg":    msg,
		"status": status,
		"value":  value,
	})
}

func getURLHandler(c *gin.Context) {
	var req request
	if err := c.ShouldBindUri(&req); err != nil {
		jsonCTX(c, http.StatusInternalServerError, "invalid url key provided", req.ID)
		return
	}

	url, err := GetURLFromUUID(db, req.ID)
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

func createURLHandler(c *gin.Context) {
	var payload struct {
		LongURL string `json:"url" binding:"required"`
	}

	if err := c.Bind(&payload); err != nil {
		jsonCTX(c, http.StatusBadRequest, "failed to bind payload url", payload.LongURL)
		return
	}

	if !validURL(payload.LongURL) {
		jsonCTX(c, http.StatusBadRequest, "invalid url string provided", payload.LongURL)
		return
	}

	url, err := GetURLFromLongURL(db, payload.LongURL)
	if err == nil {
		jsonCTX(c, http.StatusConflict, "url already exists in database", url.ID)
		return
	}

	if err := StoreURL(db, payload.LongURL); err != nil {
		jsonCTX(c, http.StatusInternalServerError, "failed to store url", payload.LongURL)
		return
	}
	url, err = GetURLFromLongURL(db, payload.LongURL)
	if err != nil {
		if err == sql.ErrNoRows {
			jsonCTX(c, http.StatusNotFound, "no result for url key provided", payload.LongURL)
			return
		}
		jsonCTX(c, http.StatusInternalServerError, err.Error(), payload.LongURL)
		return
	}

	jsonCTX(c, http.StatusOK, "successfully stored url", url.ID)
}

// IsURL check if the string is an URL.
// func IsURL(str string) bool {
// 	if str == "" || utf8.RuneCountInString(str) >= maxURLRuneCount || len(str) <= minURLRuneCount || strings.HasPrefix(str, ".") {
// 		return false
// 	}
// 	strTemp := str
// 	if strings.Contains(str, ":") && !strings.Contains(str, "://") {
// 		// support no indicated urlscheme but with colon for port number
// 		// http:// is appended so url.Parse will succeed, strTemp used so it does not impact rxURL.MatchString
// 		strTemp = "http://" + str
// 	}
// 	u, err := url.Parse(strTemp)
// 	if err != nil {
// 		return false
// 	}
// 	if strings.HasPrefix(u.Host, ".") {
// 		return false
// 	}
// 	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
// 		return false
// 	}
// 	return rxURL.MatchString(str)
// }

func validURL(str string) bool {
	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	_, err := client.Get(str)
	if err != nil {
		return false
	}
	return true
}
