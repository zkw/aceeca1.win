package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"aceeca1.win/backend/pb"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

func registerUser(e *echo.Echo, db *bbolt.DB) {
	e.GET("/ajax/user-status", func(c echo.Context) error {
		return userStatus(db, c)
	})
	e.GET("/ajax/user-login-1", func(c echo.Context) error {
		return userLogin1(db, c)
	})
	e.GET("/ajax/user-login-2", func(c echo.Context) error {
		return userLogin2(db, c)
	})
	e.GET("/ajax/user-logout", func(c echo.Context) error {
		return userLogout(db, c)
	})
	e.GET("/ajax/user-set-nick", func(c echo.Context) error {
		return userSetNick(db, c)
	})
	e.GET("/ajax/user-list", func(c echo.Context) error {
		return userList(db, c)
	})
	e.GET("/ajax/user-grant-permission-by-root", func(c echo.Context) error {
		return userGrantPermissionByRoot(db, c)
	})
}

func userStatus(db *bbolt.DB, c echo.Context) error {
	s, _ := session.Get("session", c)
	id := s.Values["id"]
	if id == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.String(http.StatusOK, getNick(db, id.(string)))
}

func userLogin1(db *bbolt.DB, c echo.Context) error {
	token := fmt.Sprintf("%08d", rand.Intn(100000000))
	db.Update(func(tx *bbolt.Tx) error {
		login := tx.Bucket([]byte(pb.Bucket_LOGIN.String()))
		login.Delete([]byte(token))
		return nil
	})
	s, _ := session.Get("session", c)
	s.Values["token"] = token
	s.Save(c.Request(), c.Response())
	return c.String(http.StatusOK, token)
}

func userLogin2(db *bbolt.DB, c echo.Context) error {
	s, _ := session.Get("session", c)
	tokenI := s.Values["token"]
	if tokenI == nil {
		return c.NoContent(http.StatusForbidden)
	}
	token := tokenI.(string)
	id := ""
	db.View(func(tx *bbolt.Tx) error {
		login := tx.Bucket([]byte(pb.Bucket_LOGIN.String()))
		result := login.Get([]byte(token))
		if result != nil {
			id = string(result)
		}
		return nil
	})
	if len(id) == 0 {
		return c.NoContent(http.StatusForbidden)
	}
	s.Values["id"] = id
	s.Save(c.Request(), c.Response())
	return c.String(http.StatusOK, getNick(db, id))
}

func userLogout(db *bbolt.DB, c echo.Context) error {
	s, _ := session.Get("session", c)
	delete(s.Values, "id")
	s.Save(c.Request(), c.Response())
	return c.NoContent(http.StatusOK)
}

func userSetNick(db *bbolt.DB, c echo.Context) error {
	s, _ := session.Get("session", c)
	id := s.Values["id"]
	if id == nil {
		return c.NoContent(http.StatusForbidden)
	}
	nick := c.QueryParam("nick")
	return c.String(http.StatusOK, setNick(db, id.(string), nick))
}

func userList(db *bbolt.DB, c echo.Context) error {
	query := c.QueryParam("query")
	result := make(map[string]string)
	db.View(func(tx *bbolt.Tx) error {
		user := tx.Bucket([]byte(pb.Bucket_USER.String()))
		user.ForEach(func(k, v []byte) error {
			vProto := pb.User{}
			proto.Unmarshal(v, &vProto)
			if strings.Contains(vProto.Nick, query) {
				result[string(k)] = vProto.Nick
			}
			return nil
		})
		return nil
	})
	return c.JSON(http.StatusOK, result)
}

func userGrantPermissionByRoot(db *bbolt.DB, c echo.Context) error {
	master := c.QueryParam("master-password")
	if !checkMasterPassword(master) {
		return c.NoContent(http.StatusForbidden)
	}
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	role := c.QueryParam("role")
	roleNumber, ok := pb.User_Permission_value[role]
	if !ok {
		return c.NoContent(http.StatusBadRequest)
	}
	rolePermission := pb.User_Permission(roleNumber)
	db.Update(func(tx *bbolt.Tx) error {
		user := tx.Bucket([]byte(pb.Bucket_USER.String()))
		result := user.Get([]byte(from))
		if result == nil {
			ok = false
			return nil
		}
		resultProto := pb.User{}
		proto.Unmarshal(result, &resultProto)
		resultProto.Role[to] = rolePermission
		data, _ := proto.Marshal(&resultProto)
		user.Put([]byte(from), data)
		return nil
	})
	if !ok {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.NoContent(http.StatusOK)
}

func isValidToken(s string) bool {
	if len(s) != 8 {
		return false
	}
	for i := 0; i < 8; i++ {
		if s[i] < '0' || '9' < s[i] {
			return false
		}
	}
	return true
}

func getNick(db *bbolt.DB, id string) string {
	nick := ""
	db.View(func(tx *bbolt.Tx) error {
		user := tx.Bucket([]byte(pb.Bucket_USER.String()))
		result := user.Get([]byte(id))
		if result != nil {
			resultProto := pb.User{}
			proto.Unmarshal(result, &resultProto)
			nick = resultProto.Nick
		}
		return nil
	})
	if len(nick) == 0 {
		return setNick(db, id, "微信用户")
	}
	return nick
}

func setNick(db *bbolt.DB, id string, nick string) string {
	if !checkNick(nick) {
		return ""
	}
	db.Update(func(tx *bbolt.Tx) error {
		user := tx.Bucket([]byte(pb.Bucket_USER.String()))
		result := user.Get([]byte(id))
		resultProto := pb.User{}
		if result != nil {
			proto.Unmarshal(result, &resultProto)
		}
		resultProto.Nick = nick
		data, _ := proto.Marshal(&resultProto)
		user.Put([]byte(id), data)
		return nil
	})
	return nick
}

func checkNick(nick string) bool {
	return 6 <= len(nick)
}
