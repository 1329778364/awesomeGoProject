package app

import (
	"net/http"
)

func (g *Gin) Response(httpCode int, errMsg string, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: httpCode,
		Msg:  errMsg,
		Data: data,
	})
	return
}

func (g *Gin) SuccessResponse(data interface{}) {
	g.C.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Msg:  "成功",
		Data: data,
	})
	return
}

func (g *Gin) BadResponse(data interface{}) {
	g.C.JSON(http.StatusBadRequest, Response{
		Code: http.StatusBadRequest,
		Msg:  "请求发生了错误",
		Data: data,
	})
	return
}

func (g *Gin) ForbiddenResponse(data interface{}) {
	g.C.JSON(http.StatusForbidden, Response{
		Code: http.StatusForbidden,
		Msg:  "请求未授权",
		Data: data,
	})
	return
}

func (g *Gin) ErrorResponse(data interface{}) {
	g.C.JSON(http.StatusInternalServerError, Response{
		Code: http.StatusInternalServerError,
		Msg:  "服务器内部错误",
		Data: data,
	})
	return
}

func (g *Gin) HasError(err error) bool {
	if err != nil {
		switch err.(type) {
		case *BadMsg:
			g.BadResponse(err.Error())
		default:
			g.ErrorResponse(err.Error())
		}
		return true
	}
	return false
}
