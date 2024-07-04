package controller

import (
	"github.com/gin-gonic/gin"
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
	splitted := strings.Split(authHeader, " ")
	tokenString := splitted[1]
	err := db.Rdb.Del(db.Rctx, tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的操作"})
	}
	c.JSON(200, gin.H{"message": "Logout successfully!"})
}
func Test(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{"username": username})
}
