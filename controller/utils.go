package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/textproto"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/charset"
	"github.com/golang-jwt/jwt"
)

const (
	jwt_secret string = "secret"
	imap_user  string = "auth@example.com"
	imap_key   string = "your_key"
)

func generatePasswdHash(pwd string, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(pwd + salt))
	bytes := hash.Sum(nil)
	//fmt.Printf("pwd:%s\nsalt:%s\nhash:%s\n", pwd, salt, hex.EncodeToString(bytes))
	return hex.EncodeToString(bytes)
}

func Sign(uid int) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  uid,
		"iat": time.Now().Unix(),
	})

	tokenString, err = token.SignedString([]byte(jwt_secret))
	return tokenString, err
}

func ParseToken(token string) (interface{}, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwt_secret), nil
	})
	if err != nil {
		return "", err
	}
	m := claim.Claims.(jwt.MapClaims)
	if uid, ok := m["id"]; ok {
		return uid, nil
	}
	return "", fmt.Errorf("fail parse")
}

func emailVerifSuccess(uid int, email string) bool {
	c, err := getIMAPClient()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()
	if _, err := c.Select("INBOX", true); err != nil {
		log.Fatal(err)
	}
	criteria := imap.NewSearchCriteria()
	header := textproto.MIMEHeader{}
	header.Add("subject", "[auth]uid:"+fmt.Sprint(uid))
	header.Add("from", email)
	criteria.Header = header
	ids, _ := c.Search(criteria)
	return len(ids) > 0
}

func getIMAPClient() (*client.Client, error) {
	imap.CharsetReader = charset.Reader

	log.Println("[IMAP] Connecting to server...")

	c, err := client.DialTLS("imap.example.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[IMAP] Connected")

	if err := c.Login(imap_user, imap_key); err != nil {
		return nil, err
	}

	log.Println("[IMAP] Logged in")

	return c, nil
}
