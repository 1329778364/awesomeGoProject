package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	_ "userSystem/docs"
	"userSystem/models"
	"userSystem/pkg/gredis"
	"userSystem/pkg/setting"
	"userSystem/pkg/validator"
	"userSystem/routers"
)

func init() {
	setting.Setup()
	models.Setup()
	gredis.Setup()
}

// @title userSystem
// @version 1.0
// @description 登录注册模块设计：密文传输+jwt身份验证
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	gin.SetMode(setting.ServerSetting.RunMode)
	validator.InitValidator()
	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info]  Set Time Zone to Asia/Chongqing")
	timeLocal, err := time.LoadLocation("Asia/Chongqing")
	if err != nil {
		log.Printf("[error] Set Time Zone to Asia/Chongqing failed %s", err)
	}
	time.Local = timeLocal
	log.Printf("[info] start server listening %s", endPoint)

	//使用https，需要将ssl.pem和ssl.key放置项目目录下
	//if errmsg := server.ListenAndServeTLS("./ssl.pem", "./ssl.key"); errmsg != nil {
	//	log.Fatalf("start https server failed %s", errmsg)
	//}
	//使用http
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("start http server failed %s", err)
	}
}
