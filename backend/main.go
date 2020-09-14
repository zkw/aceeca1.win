package main

import (
	"crypto/sha256"
	"fmt"

	"aceeca1.win/backend/pb"
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
	registerPermission(e, db)
	registerSite(e, db)
	e.Static("/", "../dist")
	e.File("*", "../dist/index.html")
	e.Logger.Fatal(e.Start(":1323"))
}

func createBuckets(db *bbolt.DB) {
	db.Update(func(tx *bbolt.Tx) error {
		for _, name := range pb.Bucket_name {
			tx.CreateBucketIfNotExists([]byte(name))
		}
		return nil
	})
}

func checkMasterPassword(master string) bool {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(master)))
	return hash == "1251145c90ac94a49f17a613c16dfd201190627bea741e4a1792a80a1a5e8dfa"
}
