package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"userSystem/pkg/app"
)

// @Summary 测试Get请求
// @Tags 测试
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/test/get [get]
func TestGet(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.SuccessResponse("测试")
}

// @Summary 测试tmpl模板文件
// @Tags 测试
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/test/tmpl [get]
func TestTmpl(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"name": "测试tmpl模板文件",
	})
}
