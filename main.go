package main

import (
	"fmt"
	"github.com/chenobdo/go-gin-example/pkg/setting"
	"github.com/chenobdo/go-gin-example/routers"
	"log"
	"net/http"
)

func main() {
	router := routers.InitRoute()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		// 请求头的最大字节数
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("Fail to listen and serve %v", err)
		return
	}
}
