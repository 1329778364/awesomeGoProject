package app

import "github.com/gin-gonic/gin"

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type BadMsg struct {
	s string
}

func NewBadMsg(text string) error {
	return &BadMsg{text}
}

func (e *BadMsg) Error() string {
	return e.s
}
