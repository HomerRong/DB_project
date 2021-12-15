package settings

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type server struct {
	RunMode  string
	Host     string
	HttpPort int
}

var ServerSetting = &server{}

type database struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

var DatabaseSetting = &database{}

type redis struct {
	Addr        string
	Password    string
	DB          int
	IdleTimeout time.Duration
}

var RedisSetting = &redis{}

var cfg *ini.File

// Setup initialize the configuration instance
func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("settings.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)

	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
