package main

import (
	"fmt"
	"github.com/chenobdo/go-gin-example/models"
	"github.com/chenobdo/go-gin-example/pkg/logging"
	"github.com/chenobdo/go-gin-example/pkg/setting"
	"github.com/chenobdo/go-gin-example/routers"
	"github.com/fvbock/endless"
	"log"
	"syscall"
)

func main() {
	// 1. 冷启动
	//router := routers.InitRoute()
	//
	//s := &http.Server{
	//	Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
	//	Handler:        router,
	//	ReadTimeout:    setting.ReadTimeout,
	//	WriteTimeout:   setting.WriteTimeout,
	//	// 请求头的最大字节数
	//	MaxHeaderBytes: 1 << 20,
	//}
	//
	//err := s.ListenAndServe()
	//if err != nil {
	//	log.Fatalf("Fail to listen and serve %v", err)
	//	return
	//}


	// 2. 热启动
	//endless.DefaultReadTimeOut = setting.ReadTimeout
	//endless.DefaultWriteTimeOut = setting.WriteTimeout
	//endless.DefaultMaxHeaderBytes = 1 << 20
	//endPoint := fmt.Sprintf(":%d", setting.HTTPPort)
	//
	//server := endless.NewServer(endPoint, routers.InitRoute())
	//server.BeforeBegin = func(add string) {
	//	log.Printf("Actual pid is %d", syscall.Getpid())
	//}
	//
	//err := server.ListenAndServe()
	//if err != nil {
	//	log.Printf("Server err: %v", err)
	//}


	// 3. shutdown 启动（守护进程）
	//router := routers.InitRoute()
	//
	//s := &http.Server{
	//	Addr: fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
	//	Handler: router,
	//	ReadTimeout: setting.ServerSetting.ReadTimeout,
	//	WriteTimeout: setting.ServerSetting.WriteTimeout,
	//	MaxHeaderBytes: 1 << 20,
	//}
	//
	//go func() {
	//	if err := s.ListenAndServe(); err != nil {
	//		log.Printf("Listen: %s\n", err)
	//	}
	//}()
	//
	//quit := make(chan os.Signal)
	//signal.Notify(quit, os.Interrupt)
	//<- quit
	//
	//log.Println("Shutdown Server ...")
	//
	//ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	//defer cancel()
	//if err := s.Shutdown(ctx); err != nil {
	//	log.Fatal("Server Shutdown:", err)
	//}
	//
	//log.Println("Server exiting")


	setting.Setup()
	models.Setup()
	logging.Setup()

	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	server := endless.NewServer(endPoint, routers.InitRoute())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
