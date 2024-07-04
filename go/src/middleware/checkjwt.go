package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"net/http"
	"shorturl/go/src/db"
	"strings"
)

var (
	jwtSecret = []byte("user_secret")
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "无效请求"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "token校验失败"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "请重新登录"})
			c.Abort()
			return
		}

		// 检查令牌是否在 Redis 中
		val, err := db.Rdb.Get(db.Rctx, tokenString).Result()
		if err == redis.Nil || val != "valid" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "请重新登录"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("username", claims["username"])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "请重新登录"})
			c.Abort()
			return
		}
		c.Next()
	}
}
