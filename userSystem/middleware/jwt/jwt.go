package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"userSystem/pkg/app"
	"userSystem/pkg/gredis"
	"userSystem/pkg/util"
	"userSystem/service/user_service"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}
		auth := c.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			//阻止挂起函数
			appG.ForbiddenResponse("请登录后再进行操作")
			c.Abort()
			return
		}
		if len(strings.Fields(auth)) > 1 {
			auth = strings.Fields(auth)[1]
		}
		// 校验token
		claims, err := util.ParseToken(auth)
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				appG.ForbiddenResponse("token 已过期")
			default:
				appG.ForbiddenResponse("token 验证失败")
			}
			c.Abort()
			return
		}
		//签发时间在 预黑名单生成时间 之前，则将其注销
		val, err := gredis.Get(user_service.LoginBlacklist(claims.Id))
		if err != nil {
			appG.ErrorResponse(err.Error())
			c.Abort()
			return
		}
		if val != "" {
			times, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				appG.ErrorResponse(err.Error())
				c.Abort()
				return
			}
			if claims.IssuedAt < times {
				appG.ForbiddenResponse("token 已过期")
				c.Abort()
				return
			}
		}
		//继续执行挂起的函数
		c.Next()
	}
}
