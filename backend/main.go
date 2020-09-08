package main

import (
	"crypto/sha1"
	"os"
	"sort"
	"strings"

	//"encoding/gob"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.etcd.io/bbolt"
)

type WxEncrypted struct {
	ToUserName string
	Encrypt    string
}

func main() {
	db, _ := bbolt.Open("database.bbolt", 0600, nil)
	db.Update(func(tx *bbolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("Article"))
		tx.CreateBucketIfNotExists([]byte("Comment"))
		tx.CreateBucketIfNotExists([]byte("Login"))
		tx.CreateBucketIfNotExists([]byte("Resource"))
		tx.CreateBucketIfNotExists([]byte("Site"))
		tx.CreateBucketIfNotExists([]byte("User"))
		return nil
	})
	defer db.Close()
	e := echo.New()
	secret := []byte(securecookie.GenerateRandomKey(32))
	e.Use(session.Middleware(sessions.NewCookieStore(secret)))
	e.GET("/ajax/wx", func(c echo.Context) error {
		return wx(db, c)
	})
	e.POST("/ajax/wx", func(c echo.Context) error {
		return wxPost(db, c)
	})
	e.GET("/ajax/user-login-1", func(c echo.Context) error {
		return userLogin1(db, c)
	})
	e.GET("/ajax/user-login-2", func(c echo.Context) error {
		return userLogin2(db, c)
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func wx(db *bbolt.DB, c echo.Context) error {
	signature := c.QueryParam("signature")
	timestamp := c.QueryParam("timestamp")
	nonce := c.QueryParam("nonce")
	echostr := c.QueryParam("echostr")
	token := os.Getenv("WX_TOKEN")
	arr := []string{token, timestamp, nonce}
	sort.Strings(arr)
	sha1Sum := fmt.Sprintf("%x", sha1.Sum([]byte(strings.Join(arr, ""))))
	if sha1Sum != signature {
		return c.NoContent(http.StatusOK)
	}
	return c.String(http.StatusOK, echostr)
}

func wxPost(db *bbolt.DB, c echo.Context) error {
	msg_signature := c.QueryParam("signature")
	nonce := c.QueryParam("nonce")
	openid := c.QueryParam("openid")
	signature := c.QueryParam("signature")
	timestamp := c.QueryParam("timestamp")
	m := WxEncrypted{}
	c.Bind(&m)
	fmt.Println(m)
	fmt.Println(c.QueryParams())
	return c.NoContent(http.StatusOK)
}

func userLogin1(db *bbolt.DB, c echo.Context) error {
	token := fmt.Sprintf("%08d", rand.Intn(100000000))
	s, _ := session.Get("session", c)
	s.Values["token"] = token
	s.Save(c.Request(), c.Response())
	return c.String(http.StatusOK, token)
}

func userLogin2(db *bbolt.DB, c echo.Context) error {
	s, _ := session.Get("session", c)
	token := s.Values["token"].(string)
	id := ""
	db.View(func(tx *bbolt.Tx) error {
		login := tx.Bucket([]byte("Login"))
		result := login.Get([]byte(token))
		if result != nil {
			id = string(result)
		}
		return nil
	})
	if len(id) == 0 {
		return c.NoContent(http.StatusOK)
	}
	s.Values["id"] = id
	s.Save(c.Request(), c.Response())
	return c.String(http.StatusOK, id)
}
