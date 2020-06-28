package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	Name     string
}

type Redis struct {
	Host     string
	Password string
}

type Mail struct {
	User string
	Pass string
	Host string
	Port string
}

var ServerSetting = &Server{}

var DatabaseSetting = &Database{}

var RedisSetting = &Redis{}

var MailSetting = &Mail{}

var cfg *ini.File

// 程序初始化配置
func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}
	mapTo("server", ServerSetting)
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)
	mapTo("mail", MailSetting)
}

// 在 go-ini 中可以采用 MapTo 的方式来映射结构体
// 读取conf/app.ini的section信息，映射到结构体中
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s errmsg: %v", section, err)
	}
}
