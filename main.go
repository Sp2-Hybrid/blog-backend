
package main

import (
	"fmt"
	"github.com/fvbock/endless"
	"go-gin-example/models"
	"go-gin-example/pkg/logging"
	"go-gin-example/pkg/setting"
	"go-gin-example/routers"
	"log"
	"syscall"
)

func main() {
	setting.Setup()
	fmt.Println("-------------------setting.Setup over--------------------")
	models.Setup()
	fmt.Println("-------------------models.Setup over--------------------")
	logging.Setup()
	fmt.Println("-------------------logging.Setup over--------------------")
	router := routers.InitRouter()

	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	server := endless.NewServer(endPoint, router)
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err!=nil{
		log.Printf("Server err: %v", err)
	}
}