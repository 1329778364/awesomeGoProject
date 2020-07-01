package jwt

import (
	"github.com/gin-gonic/gin"
	"userSystem/pkg/app"
	"userSystem/pkg/util"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}
		auth := c.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			//阻止挂起函数
			c.Abort()
			appG.ForbiddenResponse("认证失败，请重新登录")
			return
		}
		// 校验token
		_, err := util.ParseToken(auth)
		if err != nil {
			//阻止挂起函数
			c.Abort()
			appG.ForbiddenResponse("token 验证失败")
			return
		}
		//继续执行挂起的函数
		c.Next()
	}
}
