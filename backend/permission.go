package main

import (
	"net/http"

	"aceeca1.win/backend/pb"
	"github.com/labstack/echo/v4"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

func registerPermission(e *echo.Echo, db *bbolt.DB) {
	e.POST("/ajax/permission-view", func(c echo.Context) error {
		return permissionView(db, c)
	})
	e.POST("/ajax/permission-edit", func(c echo.Context) error {
		return permissionEdit(db, c)
	})
}

type requestPermission struct {
	MasterPassword  string
	Name            string
	PermissionProto string
}

func permissionView(db *bbolt.DB, c echo.Context) error {
	r := requestPermission{}
	c.Bind(&r)
	if !checkMasterPassword(r.MasterPassword) {
		return c.String(http.StatusForbidden, "主密码错误")
	}
	var resultText string
	db.View(func(tx *bbolt.Tx) error {
		permission := tx.Bucket([]byte(pb.Bucket_PERMISSION.String()))
		result := permission.Get([]byte(r.Name))
		if result != nil {
			resultProto := pb.Permission{}
			proto.Unmarshal(result, &resultProto)
			resultText = prototext.Format(&resultProto)
		}
		return nil
	})
	return c.String(http.StatusOK, resultText)
}

func permissionEdit(db *bbolt.DB, c echo.Context) error {
	r := requestPermission{}
	c.Bind(&r)
	if !checkMasterPassword(r.MasterPassword) {
		return c.String(http.StatusForbidden, "主密码错误")
	}
	permissionProto := pb.Permission{}
	err := prototext.Unmarshal([]byte(r.PermissionProto), &permissionProto)
	if err != nil {
		return c.String(http.StatusBadRequest, "语法错误")
	}
	db.Update(func(tx *bbolt.Tx) error {
		permission := tx.Bucket([]byte(pb.Bucket_PERMISSION.String()))
		if len(permissionProto.User) == 0 {
			permission.Delete([]byte(r.Name))
		} else {
			data, _ := proto.Marshal(&permissionProto)
			permission.Put([]byte(r.Name), data)
		}
		return nil
	})
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
