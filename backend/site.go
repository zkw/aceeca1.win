package main

import (
	"net/http"
	"strings"

	"aceeca1.win/backend/pb"
	"github.com/labstack/echo/v4"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

type requestSite struct {
	MasterPassword string
	ID             string
	SiteProto      string
}

func registerSite(e *echo.Echo, db *bbolt.DB) {
	e.GET("/ajax/site-information", func(c echo.Context) error {
		return siteInformation(db, c)
	})
	e.POST("/ajax/site-view", func(c echo.Context) error {
		return siteView(db, c)
	})
	e.POST("/ajax/site-edit", func(c echo.Context) error {
		return siteEdit(db, c)
	})
}

func siteInformation(db *bbolt.DB, c echo.Context) error {
	host := c.Request().Host
	sub := strings.SplitN(host, ".", 2)
	answer := make(map[string]string)
	db.View(func(tx *bbolt.Tx) error {
		site := tx.Bucket([]byte(pb.Bucket_SITE.String()))
		result := site.Get([]byte([]byte(sub[0])))
		if result != nil {
			resultProto := pb.Site{}
			proto.Unmarshal(result, &resultProto)
			documentSide := resultProto.AllowedUser
			userSide := permissionForUser(tx, c)
			if permissionMatch(userSide, documentSide) <= pb.DocumentPermission_READER {
				answer["Name"] = resultProto.Name
			}
		}
		return nil
	})
	return c.JSON(http.StatusOK, answer)
}

func siteView(db *bbolt.DB, c echo.Context) error {
	r := requestSite{}
	c.Bind(&r)
	if !checkMasterPassword(r.MasterPassword) {
		return c.String(http.StatusForbidden, "主密码错误")
	}
	var resultText string
	db.View(func(tx *bbolt.Tx) error {
		site := tx.Bucket([]byte(pb.Bucket_SITE.String()))
		result := site.Get([]byte(r.ID))
		if result != nil {
			resultProto := pb.Site{}
			proto.Unmarshal(result, &resultProto)
			resultText = prototext.Format(&resultProto)
		}
		return nil
	})
	return c.String(http.StatusOK, resultText)
}

func siteEdit(db *bbolt.DB, c echo.Context) error {
	r := requestSite{}
	c.Bind(&r)
	if !checkMasterPassword(r.MasterPassword) {
		return c.String(http.StatusForbidden, "主密码错误")
	}
	siteProto := pb.Site{}
	err := prototext.Unmarshal([]byte(r.SiteProto), &siteProto)
	if err != nil {
		return c.String(http.StatusBadRequest, "语法错误")
	}
	db.Update(func(tx *bbolt.Tx) error {
		site := tx.Bucket([]byte(pb.Bucket_SITE.String()))
		if len(siteProto.Name) == 0 {
			site.Delete([]byte(r.ID))
		} else {
			data, _ := proto.Marshal(&siteProto)
			site.Put([]byte(r.ID), data)
		}
		return nil
	})
	return c.NoContent(http.StatusOK)
}
