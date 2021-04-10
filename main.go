package main

import (
	"fmt"
	"go-gin-example/pkg/setting"
	"go-gin-example/routers"
	"net/http"
)

func main() {
	router := routers.InitRouter()

	s := &http.Server{
		WriteTimeout:setting.WriteTimeout,
		Addr:fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:router,
		ReadTimeout:setting.ReadTimeout,
		MaxHeaderBytes:1<<20,
		}
	s.ListenAndServe()
}
