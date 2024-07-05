package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"net/url"
	"shorturl/go/src/db"
	"shorturl/go/src/tool"
	"strconv"
	"time"
)

func CreateShortUrl(c *gin.Context) {
	longUrl := c.PostForm("url")
	customShort := c.PostForm("customShort")
	if len(longUrl) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "url不能为空"})
	}

	//限制单个用户调用此api次数，30分钟内最多生产五个短连接
	username, _ := c.Get("username")
	usernameStr := username.(string)
	val, err := db.Rdb.Get(db.Rctx, usernameStr).Result()
	if err == redis.Nil {
		_ = db.Rdb.Set(db.Rctx, usernameStr, 5, 30*60*time.Second).Err()
	} else {
		val, _ = db.Rdb.Get(db.Rctx, usernameStr).Result()
		limit, _ := strconv.Atoi(val)
		if limit <= 0 {
			ext, _ := db.Rdb.TTL(db.Rctx, usernameStr).Result()
			c.JSON(http.StatusBadRequest, gin.H{"message": "已超出可创建个数", "retry_time": ext / time.Second / time.Minute})
			return
		}
	}
	//校验传输的原url是否合理
	_, err = url.ParseRequestURI(longUrl)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的连接"})
		return
	}
	//校验是否为服务端口
	if !tool.RemoveDoMain(longUrl) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的连接"})
	}
	//强制执行ssl证书
	longUrl = tool.EnforceHttp(longUrl)
	//校验自定义短连接是否存在
	var short_id string
	if customShort == "" {
		short_id = uuid.New().String()[:6]
	} else {
		short_id = customShort
	}
	val, _ = db.Rdb.Get(db.Rctx, short_id).Result()
	if val != "" {
		c.JSON(http.StatusForbidden, gin.H{"message": "该链接已被使用"})
	}

	err = db.Rdb.Set(db.Rctx, short_id, longUrl, 24*3600*time.Second).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无法连接到数据服务"})
		return
	}

	db.Rdb.Decr(db.Rctx, usernameStr)
	less, _ := db.Rdb.Get(db.Rctx, usernameStr).Result()
	ttl, _ := db.Rdb.Get(db.Rctx, short_id).Result()
	ttlTime, _ := time.ParseDuration(ttl)
	c.JSON(http.StatusOK, gin.H{"url": "http://localhost:9900/api/" + short_id, "message": "你还能创建" + less + "次", "retry_time": ttlTime})
}

func ResolveShortUrl(c *gin.Context) {
	shortUrl := c.Param("shortUrl")

	value, err := db.Rdb.Get(db.Rctx, shortUrl).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "没有找到短连接"})
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "内部服务器错误"})
	}
	c.Redirect(302, value)
}
