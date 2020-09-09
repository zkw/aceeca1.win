package main

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.etcd.io/bbolt"
)

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
	registerWX(e, db)
	registerUser(e, db)
	e.Logger.Fatal(e.Start(":1323"))
}
