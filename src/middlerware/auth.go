package middlerware

import (
	"github.com/chenxuan520/qiniuserver/config"
	"github.com/gin-gonic/gin"
)

func PasswdAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.GlobalConfig.Host.Password == "" {
			c.Next()
			return
		}

		// 获取请求头中的password
		password := c.GetHeader("Authorization")
		if password == "" {
			c.JSON(401, gin.H{
				"code": 401,
				"msg":  "password is empty",
			})
			c.Abort()
			return
		}

		// 验证token是否正确
		if password != config.GlobalConfig.Host.Password {
			c.JSON(401, gin.H{
				"code": 401,
				"msg":  "password is error",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
