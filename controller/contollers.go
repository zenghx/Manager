package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zenghx/manager/db"
)

func SignUp(c *gin.Context) {
	info := &db.UserInfo{}
	c.BindJSON(info)
	if db.IsNewEmail(info.Email) {
		salt := fmt.Sprintf("%x", time.Now().UnixNano())
		pwdhash := generatePasswdHash(info.Pwd, salt)
		//fmt.Println(db.User{Email: info.Email, Salt: salt, PwdHash: pwdhash})
		db.CreateUser(&db.User{Email: info.Email, Salt: salt, PwdHash: pwdhash})
		c.JSON(http.StatusCreated, gin.H{
			"code": 200,
			"msg":  "注册成功",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 409,
			"msg":  "电子邮箱地址已被注册",
		})
	}
}

func SignIn(c *gin.Context) {
	info := &db.UserInfo{}
	c.BindJSON(info)
	user := db.GetUserByEmail(info.Email)
	if user.UID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "电子邮箱地址或密码不正确。",
		})
	} else if user.PwdHash == generatePasswdHash(info.Pwd, user.Salt) {
		cookieStr, err := Sign(user.UID)
		if err != nil {
			log.Fatal(err)
		}
		c.SetCookie("token", cookieStr, 3600*24*3650, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 403,
			"msg":  "电子邮箱地址或密码不正确。",
		})
	}
}

func Auth(c *gin.Context, r *gin.Engine) {
	_, err := c.Request.Cookie("token")
	if err != nil {
		c.Abort()
		c.Request.URL.Path = "/"
		r.HandleContext(c)
	} else {
		c.Next()
	}
}

func RefreshEmail(c *gin.Context) {
	_, err := c.Request.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "请先登录",
		})
	} else {
		u := db.User{}
		token, _ := c.Request.Cookie("token")
		if id, err := ParseToken(token.Value); err == nil {
			u = db.GetUserById(int(id.(float64)))
		} else {
			fmt.Println(err)
		}
		if emailVerifSuccess(u.UID, u.Email) {
			db.UpdateVerif(u.UID)
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "电子邮件地址验证成功！",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 202,
				"msg":  "请耐心等待，若时间过久请尝试重新发送邮件。",
			})
		}
	}
}
