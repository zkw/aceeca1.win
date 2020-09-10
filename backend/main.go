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
	return hash == "bda88f8614f72aaecdd0864f812a88493ed9ed093d3ad5fcbedd8c0f57b84baa"
}
