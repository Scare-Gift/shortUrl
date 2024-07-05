package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"shorturl/go/src/db"
	"shorturl/go/src/service"
	"shorturl/go/src/tool"
	"strings"
	"time"
)

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户名或密码不能为空"})
		return
	}
	user, err := service.GetUserInfo(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "发生了未知错误"})
		return
	}
	if username != user.Username && password != user.Password {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户名或密码错误"})
		return
	}
	tokenString, _ := tool.GenerateToken(username)
	db.Rdb.Set(db.Rctx, tokenString, "valid", time.Hour*24)
	c.JSON(http.StatusOK, gin.H{"token": tokenString, "message": "登录成功"})
}

func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	split := strings.Split(authHeader, " ")
	tokenString := split[1]
	val, _ := db.Rdb.Exists(db.Rctx, tokenString).Result()
	log.Print(val)
	if val == 1 {
		db.Rdb.Del(db.Rctx, tokenString)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "非法token"})
		return
	}
	c.JSON(200, gin.H{"message": "Logout successfully!"})
}
