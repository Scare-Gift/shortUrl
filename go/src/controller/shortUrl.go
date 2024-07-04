package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateShortUrl(c *gin.Context) {
	url := c.PostForm("url")
	if len(url) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "url不能为空"})
	}

	//限制单个用户调用此api次数，30分钟内最多生产五个短连接

	//校验url是否合理

	//校验是否为服务端口

}
