package main

import (
	"aceeca1.win/backend/proto"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.etcd.io/bbolt"
)

func main() {
	db, _ := bbolt.Open("database.bbolt", 0600, nil)
	defer db.Close()
	createBuckets(db)
	e := echo.New()
	secret := []byte(securecookie.GenerateRandomKey(32))
	e.Use(session.Middleware(sessions.NewCookieStore(secret)))
	registerWX(e, db)
	registerUser(e, db)
	e.Logger.Fatal(e.Start(":1323"))
}

func createBuckets(db *bbolt.DB) {
	db.Update(func(tx *bbolt.Tx) error {
		for _, name := range proto.Bucket_name {
			tx.CreateBucketIfNotExists([]byte(name))
		}
		return nil
	})
}
