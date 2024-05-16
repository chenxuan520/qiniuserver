package middlerware

import (
	"time"

	"github.com/chenxuan520/qiniuserver/config"
	"github.com/gin-gonic/gin"
)

func PasswdAuth() gin.HandlerFunc {
	var nextAllowTryPwdTime int64 = 0
	return func(c *gin.Context) {
		if config.GlobalConfig.Host.Password == "" {
			c.Next()
			return
		}
		// 判断是否到达时间,避免暴力攻击
		if time.Now().Unix() < nextAllowTryPwdTime {
			// 错误时间加3s
			nextAllowTryPwdTime = time.Now().Unix() + 3

			c.JSON(401, gin.H{
				"code": 401,
				"msg":  "password is too many error,please try again later",
			})
			c.Abort()
			return
		}

		// 获取请求头中的password
		password := c.GetHeader("Authorization")

		// 验证token是否正确
		if password != config.GlobalConfig.Host.Password {
			// 错误时间加3s
			nextAllowTryPwdTime = time.Now().Unix() + 3

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
