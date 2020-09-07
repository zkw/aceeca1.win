package main

import (
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

func main() {
	db, _ := bbolt.Open("database.bbolt", 0600, nil)
	db.Update(func(tx *bbolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("User"))
		tx.CreateBucketIfNotExists([]byte("Article"))
		tx.CreateBucketIfNotExists([]byte("Comment"))
		tx.CreateBucketIfNotExists([]byte("Site"))
		tx.CreateBucketIfNotExists([]byte("Resource"))
		return nil
	})
	defer db.Close()
	e := echo.New()
	secret := []byte(securecookie.GenerateRandomKey(32))
	e.Use(session.Middleware(sessions.NewCookieStore(secret)))
	e.GET("/ajax/user-login-1", func(c echo.Context) error {
		return userLogin1(db, c)
	})
	e.GET("/ajax/user-login-2", func(c echo.Context) error {
		return userLogin2(db, c)
	})
	e.Logger.Fatal(e.Start(":1323"))
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
	token := s.Values["token"]
	db.View(func(tx *bbolt.Tx) error {
		return nil
	})
}
