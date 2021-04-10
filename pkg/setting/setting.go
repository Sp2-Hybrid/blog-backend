package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

var (
	Cfg *ini.File
	RunMode string
	HTTPPort int
	ReadTimeout time.Duration
	WriteTimeout time.Duration

	PageSize int
)

func init()  {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err!=nil{
		log.Fatalf("Fail to parse 'conf/app.ini' : %v", err)
	}
	LoadBase()
	LoadApp()
	LoadServer()
}

func LoadBase()  {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func LoadServer()  {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout =  time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

//func LoadDatabase(){
//	sec, err := Cfg.GetSection("database")
//	if err!=nil{
//		log.Fatalf("Fail to get section 'database' : %v", err)
//	}
//
//}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err!=nil{
		log.Fatalf("Fail to get section 'app' : %v", err)
	}
	PageSize = sec.Key("PAGE_SIZE").MustInt()
}
