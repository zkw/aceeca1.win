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

	"github.com/labstack/echo/v4"
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
}

const wxResponseTemplate = `
<xml>
	<Encrypt><![CDATA[%s]]></Encrypt>
	<MsgSignature><![CDATA[%s]]></MsgSignature>
	<TimeStamp>%s</TimeStamp>
	<Nonce><![CDATA[%s]]></Nonce>
</xml>
`

func registerWX(e *echo.Echo, db *bbolt.DB) {
	e.GET("/ajax/wx", func(c echo.Context) error {
		return wx(db, c)
	})
	e.POST("/ajax/wx", func(c echo.Context) error {
		return wxPost(db, c)
	})
}

func wx(db *bbolt.DB, c echo.Context) error {
	signature := c.QueryParam("signature")
	timestamp := c.QueryParam("timestamp")
	nonce := c.QueryParam("nonce")
	echostr := c.QueryParam("echostr")
	if wxSign(timestamp, nonce, "") != signature {
		return c.NoContent(http.StatusOK)
	}
	return c.String(http.StatusOK, echostr)
}

func wxPost(db *bbolt.DB, c echo.Context) error {
	msg_signature := c.QueryParam("msg_signature")
	nonce := c.QueryParam("nonce")
	openid := c.QueryParam("openid")
	signature := c.QueryParam("signature")
	timestamp := c.QueryParam("timestamp")
	encrypted := WxEncrypted{}
	c.Bind(&encrypted)
	if wxSign(timestamp, nonce, "") != signature {
		return c.NoContent(http.StatusOK)
	}
	if wxSign(timestamp, nonce, encrypted.Encrypt) != msg_signature {
		return c.NoContent(http.StatusOK)
	}
	decryptedBytes := wxDecrypt(encrypted.Encrypt)
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
	decrypted.CreateTime = int(time.Now().Unix())
	decryptedBytes, _ = xml.Marshal(decrypted)
	randomString := fmt.Sprintf("%016x", rand.Uint64())
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(decryptedBytes)))
	decryptedBytes = bytes.Join([][]byte{
		[]byte(randomString),
		lenBytes,
		decryptedBytes,
		[]byte(os.Getenv("WX_APPID")),
	}, []byte{})
	encrypted.Encrypt = wxEncrypt(decryptedBytes)
	nonce = fmt.Sprintf("%d", rand.Int31())
	timestamp = fmt.Sprintf("%d", time.Now().Unix())
	msg_signature = wxSign(timestamp, nonce, encrypted.Encrypt)
	result := fmt.Sprintf(wxResponseTemplate,
		encrypted.Encrypt, msg_signature, timestamp, nonce)
	return c.String(http.StatusOK, result)
}

func wxSign(timestamp, nonce, message string) string {
	token := os.Getenv("WX_TOKEN")
	arr := []string{token, timestamp, nonce, message}
	sort.Strings(arr)
	return fmt.Sprintf("%x", sha1.Sum([]byte(strings.Join(arr, ""))))
}

func wxDecrypt(encrypted string) []byte {
	key, _ := base64.StdEncoding.DecodeString(os.Getenv("WX_KEY") + "=")
	block, _ := aes.NewCipher(key)
	decrypter := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	byteArray, _ := base64.StdEncoding.DecodeString(encrypted)
	decrypter.CryptBlocks(byteArray, byteArray)
	return byteArray[:len(byteArray)-int(byteArray[len(byteArray)-1])]
}

func wxEncrypt(decrypted []byte) string {
	key, _ := base64.StdEncoding.DecodeString(os.Getenv("WX_KEY") + "=")
	block, _ := aes.NewCipher(key)
	encrypter := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	amountPad := block.BlockSize() - len(decrypted)%block.BlockSize()
	byteArray := append(decrypted,
		bytes.Repeat([]byte{byte(amountPad)}, amountPad)...)
	encrypter.CryptBlocks(byteArray, byteArray)
	return base64.StdEncoding.EncodeToString(byteArray)
}
