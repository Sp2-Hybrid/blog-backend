
package main

import (
	"context"
	"fmt"
	"go-gin-example/pkg/setting"
	"go-gin-example/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,

		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := s.ListenAndServe(); err != nil {
			log.Printf("Listen: %s\n", err)
		}
	}()

	// 创建一个接收信号的通道
	quit := make(chan os.Signal)
	// kill和kill -1会默认发送syscall.SIGTERM信号
	// kill -2 会发送syscall.SIGINT 信号
	// kill -9 发送syscall.SIGKILL 信号，我们常见的Ctrl+C就是触发这个信号
	// signal.Notify 将收到的信号转发给quit，这一步并不会阻塞
	signal.Notify(quit, os.Interrupt)

	// 在这一步会阻塞，当接收到信号才会往下执行
	<- quit

	log.Println("Shutdown Server ...")
	// 创建一个超时5秒的context
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	// 5秒内关闭服务（将未处理完的请求处理完再关闭服务），超时5秒就自动退出
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}