package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"aceeca1.win/backend/pb"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

type requestUserByMaster struct {
	MasterPassword string
	ID             string
	UserProto      string
}

type requestUserByAdministrator struct {
	User       string
	Permission string
	Role       string
}

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
	e.GET("/ajax/user-role-as-administrator", func(c echo.Context) error {
		return userRoleAsAdministrator(db, c)
	})
	e.POST("/ajax/user-edit-permission-by-administrator", func(c echo.Context) error {
		return userEditPermissionByAdministrator(db, c)
	})
	e.POST("/ajax/user-remove-permission-by-administrator", func(c echo.Context) error {
		return userRemovePermissionByAdministrator(db, c)
	})
	e.POST("/ajax/user-view-permission-by-root", func(c echo.Context) error {
		return userViewPermissionByRoot(db, c)
	})
	e.POST("/ajax/user-edit-permission-by-root", func(c echo.Context) error {
		return userEditPermissionByRoot(db, c)
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
	if setNick(db, id.(string), nick) != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.NoContent(http.StatusOK)
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

func userRoleAsAdministrator(db *bbolt.DB, c echo.Context) error {
	s, _ := session.Get("session", c)
	id := s.Values["id"]
	if id == nil {
		return c.String(http.StatusForbidden, "用户没有登录")
	}
	resultArray := []string{}
	db.View(func(tx *bbolt.Tx) error {
		user := tx.Bucket([]byte(pb.Bucket_USER.String()))
		result := user.Get([]byte(id.(string)))
		resultProto := pb.User{}
		proto.Unmarshal(result, &resultProto)
		for k, v := range resultProto.Role {
			if v == pb.User_ADMINISTRATOR {
				resultArray = append(resultArray, k)
			}
		}
		return nil
	})
	return c.JSON(http.StatusOK, resultArray)
}

func userEditPermissionByAdministrator(db *bbolt.DB, c echo.Context) error {
	s, _ := session.Get("session", c)
	id := s.Values["id"]
	if id == nil {
		return c.String(http.StatusForbidden, "用户没有登录")
	}
	r := requestUserByAdministrator{}
	c.Bind(&r)
	roleNumber, ok := pb.User_Permission_value[r.Role]
	if !ok {
		return c.String(http.StatusBadRequest, "语法错误")
	}
	role := pb.User_Permission(roleNumber)
	var err error
	db.Update(func(tx *bbolt.Tx) error {
		user := tx.Bucket([]byte(pb.Bucket_USER.String()))
		me := user.Get([]byte(id.(string)))
		meProto := pb.User{}
		proto.Unmarshal(me, &meProto)
		if meProto.Role == nil {
			err = errors.New("用户没有任何权限")
			return nil
		}
		mePermission, ok := meProto.Role[r.Permission]
		if !ok {
			err = errors.New("用户没有对应权限")
			return nil
		}
		if mePermission != pb.User_ADMINISTRATOR {
			err = errors.New("用户不是此权限的管理员")
			return nil
		}
		result := user.Get([]byte(r.User))
		resultProto := pb.User{}
		proto.Unmarshal(result, &resultProto)
		resultProto.Role[r.Permission] = role
		data, _ := proto.Marshal(&resultProto)
		user.Put([]byte(r.User), data)
		return nil
	})
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func userRemovePermissionByAdministrator(db *bbolt.DB, c echo.Context) error {
	s, _ := session.Get("session", c)
	id := s.Values["id"]
	if id == nil {
		return c.String(http.StatusForbidden, "用户没有登录")
	}
	r := requestUserByAdministrator{}
	c.Bind(&r)
	var err error
	db.Update(func(tx *bbolt.Tx) error {
		user := tx.Bucket([]byte(pb.Bucket_USER.String()))
		me := user.Get([]byte(id.(string)))
		meProto := pb.User{}
		proto.Unmarshal(me, &meProto)
		if meProto.Role == nil {
			err = errors.New("用户没有任何权限")
			return nil
		}
		mePermission, ok := meProto.Role[r.Permission]
		if !ok {
			err = errors.New("用户没有对应权限")
			return nil
		}
		if mePermission != pb.User_ADMINISTRATOR {
			err = errors.New("用户不是此权限的管理员")
			return nil
		}
		result := user.Get([]byte(r.User))
		resultProto := pb.User{}
		proto.Unmarshal(result, &resultProto)
		delete(resultProto.Role, r.Permission)
		data, _ := proto.Marshal(&resultProto)
		user.Put([]byte(r.User), data)
		return nil
	})
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func userViewPermissionByRoot(db *bbolt.DB, c echo.Context) error {
	r := requestUserByMaster{}
	c.Bind(&r)
	if !checkMasterPassword(r.MasterPassword) {
		return c.String(http.StatusForbidden, "主密码错误")
	}
	var resultText string
	db.View(func(tx *bbolt.Tx) error {
		user := tx.Bucket([]byte(pb.Bucket_USER.String()))
		result := user.Get([]byte(r.ID))
		if result != nil {
			resultProto := pb.User{}
			proto.Unmarshal(result, &resultProto)
			ruleProto := pb.User{}
			ruleProto.Role = resultProto.Role
			resultText = prototext.Format(&ruleProto)
		}
		return nil
	})
	if len(resultText) == 0 {
		return c.String(http.StatusBadRequest, "无此用户")
	}
	return c.String(http.StatusOK, resultText)
}

func userEditPermissionByRoot(db *bbolt.DB, c echo.Context) error {
	r := requestUserByMaster{}
	c.Bind(&r)
	if !checkMasterPassword(r.MasterPassword) {
		return c.String(http.StatusForbidden, "主密码错误")
	}
	userProto := pb.User{}
	err := prototext.Unmarshal([]byte(r.UserProto), &userProto)
	if err != nil {
		return c.String(http.StatusBadRequest, "语法错误")
	}
	db.Update(func(tx *bbolt.Tx) error {
		user := tx.Bucket([]byte(pb.Bucket_USER.String()))
		result := user.Get([]byte(r.ID))
		if result == nil {
			err = errors.New("无此用户")
			return nil
		}
		resultProto := pb.User{}
		proto.Unmarshal(result, &resultProto)
		resultProto.Role = userProto.Role
		data, _ := proto.Marshal(&resultProto)
		user.Put([]byte(r.ID), data)
		return nil
	})
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
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
		defaultNick := "微信用户"
		setNick(db, id, defaultNick)
		return defaultNick
	}
	return nick
}

func setNick(db *bbolt.DB, id string, nick string) error {
	err := checkNick(nick)
	if err != nil {
		return err
	}
	return db.Update(func(tx *bbolt.Tx) error {
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
}

func checkNick(nick string) error {
	if len(nick) < 6 {
		return errors.New("昵称太短")
	}
	return nil
}
