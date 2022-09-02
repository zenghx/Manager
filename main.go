package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zenghx/manager/controller"
	"github.com/zenghx/manager/db"
)

func main() {
	db.InitDBConn()
	r := gin.Default()
	r.LoadHTMLGlob("tmpl/html/*")
	r.StaticFile("/initn.js", "./tmpl/static/initn.js")
	r.StaticFile("/favicon.ico", "./tmpl/static/favicon.ico")
	r.GET("/", func(ctx *gin.Context) {
		if _, err := ctx.Request.Cookie("token"); err == nil {
			ctx.Request.URL.Path = "/user"
			r.HandleContext(ctx)
		} else {
			ctx.HTML(http.StatusOK, "index.html", gin.H{})
		}
	})
	r.GET("/user", func(ctx *gin.Context) {
		u := db.User{}
		token, err := ctx.Request.Cookie("token")
		if err != nil {
			ctx.Request.URL.Path = "/"
			r.HandleContext(ctx)
		} else {
			if id, err := controller.ParseToken(token.Value); err == nil {
				u = db.GetUserById(int(id.(float64)))
			} else {
				fmt.Println(err)
			}
			ctx.HTML(http.StatusOK, "user.html", gin.H{
				"uid":        u.UID,
				"email":      u.Email,
				"uuid":       string(u.UUID),
				"verified":   u.Verified,
				"authorized": u.Authorized,
			})
		}

	})
	r.POST("/signup", controller.SignUp)
	r.POST("/signin", controller.SignIn)
	r.GET("/check", controller.RefreshEmail)
	r.Run(":2291")
}
