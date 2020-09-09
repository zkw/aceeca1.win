package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.etcd.io/bbolt"
)

type WxEncrypted struct {
	ToUserName string
	Encrypt    string
}

type WxDecrypted struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int
	MsgType      string
	Content      string
	MsgId        int64
}

const wxResponseTemplate = `
<xml>
	<Encrypt><![CDATA[%s]]></Encrypt>
	<MsgSignature><![CDATA[%s]]></MsgSignature>
	<TimeStamp>%s</TimeStamp>
	<Nonce><![CDATA[%s]]></Nonce>
</xml>
`

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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080"},
	}))
	e.GET("/ajax/wx", func(c echo.Context) error {
		return wx(db, c)
	})
	e.POST("/ajax/wx", func(c echo.Context) error {
		return wxPost(db, c)
	})
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
	e.Logger.Fatal(e.Start(":1323"))
}

func sign(timestamp, nonce, message string) string {
	token := os.Getenv("WX_TOKEN")
	arr := []string{token, timestamp, nonce, message}
	sort.Strings(arr)
	return fmt.Sprintf("%x", sha1.Sum([]byte(strings.Join(arr, ""))))
}

func decrypt(encrypted string) []byte {
	key, _ := base64.StdEncoding.DecodeString(os.Getenv("WX_KEY") + "=")
	block, _ := aes.NewCipher(key)
	decrypter := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	byteArray, _ := base64.StdEncoding.DecodeString(encrypted)
	decrypter.CryptBlocks(byteArray, byteArray)
	return byteArray[:len(byteArray)-int(byteArray[len(byteArray)-1])]
}

func encrypt(decrypted []byte) string {
	key, _ := base64.StdEncoding.DecodeString(os.Getenv("WX_KEY") + "=")
	block, _ := aes.NewCipher(key)
	encrypter := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	amountPad := block.BlockSize() - len(decrypted)%block.BlockSize()
	byteArray := append(decrypted,
		bytes.Repeat([]byte{byte(amountPad)}, amountPad)...)
	encrypter.CryptBlocks(byteArray, byteArray)
	return base64.StdEncoding.EncodeToString(byteArray)
}

func wx(db *bbolt.DB, c echo.Context) error {
	signature := c.QueryParam("signature")
	timestamp := c.QueryParam("timestamp")
	nonce := c.QueryParam("nonce")
	echostr := c.QueryParam("echostr")
	if sign(timestamp, nonce, "") != signature {
		return c.NoContent(http.StatusOK)
	}
	return c.String(http.StatusOK, echostr)
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

func wxPost(db *bbolt.DB, c echo.Context) error {
	msg_signature := c.QueryParam("msg_signature")
	nonce := c.QueryParam("nonce")
	openid := c.QueryParam("openid")
	signature := c.QueryParam("signature")
	timestamp := c.QueryParam("timestamp")
	encrypted := WxEncrypted{}
	c.Bind(&encrypted)
	if sign(timestamp, nonce, "") != signature {
		return c.NoContent(http.StatusOK)
	}
	if sign(timestamp, nonce, encrypted.Encrypt) != msg_signature {
		return c.NoContent(http.StatusOK)
	}
	decryptedBytes := decrypt(encrypted.Encrypt)
	decrypted := WxDecrypted{}
	xml.Unmarshal(decryptedBytes[20:len(decryptedBytes)-18], &decrypted)
	if !isValidToken(decrypted.Content) {
		return c.NoContent(http.StatusOK)
	}
	db.Update(func(tx *bbolt.Tx) error {
		tx.Bucket([]byte("Login")).Put([]byte(decrypted.Content), []byte(openid))
		return nil
	})
	decrypted.FromUserName = decrypted.ToUserName
	decrypted.ToUserName = openid
	decrypted.MsgType = "text"
	decrypted.Content = "已收到验证码，请点击网页上的按钮登录。"
	decrypted.MsgId = rand.Int63()
	decrypted.CreateTime = int(time.Now().Unix())
	decryptedBytes, _ = xml.Marshal(decrypted)
	fmt.Println(string(decryptedBytes))
	randomString := fmt.Sprintf("%x", rand.Uint64())
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(decryptedBytes)))
	decryptedBytes = bytes.Join([][]byte{
		[]byte(randomString),
		lenBytes,
		decryptedBytes,
		[]byte(os.Getenv("WX_APPID")),
	}, []byte{})
	encrypted.Encrypt = encrypt(decryptedBytes)
	nonce = fmt.Sprintf("%d", rand.Int31())
	timestamp = fmt.Sprintf("%d", time.Now().Unix())
	msg_signature = sign(timestamp, nonce, encrypted.Encrypt)
	result := fmt.Sprintf(wxResponseTemplate,
		encrypted.Encrypt, msg_signature, timestamp, nonce)
	fmt.Println(result)
	return c.String(http.StatusOK, result)
}

func userStatus(db *bbolt.DB, c echo.Context) error {
	s, _ := session.Get("session", c)
	id := s.Values["id"]
	if id == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.String(http.StatusOK, id.(string))
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
	tokenI := s.Values["token"]
	if tokenI == nil {
		return c.NoContent(http.StatusOK)
	}
	token := tokenI.(string)
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

func userLogout(db *bbolt.DB, c echo.Context) error {
	s, _ := session.Get("session", c)
	delete(s.Values, "id")
	s.Save(c.Request(), c.Response())
	return c.NoContent(http.StatusOK)
}
