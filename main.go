package main

import (
	"bytes"
	"gin/config"
	"gin/models"
	"gin/routes"
	"gin/ws"
	"github.com/spf13/viper"
	"log"
)

func init() {
	viper.SetConfigType("json") // 设置配置文件的类型
	if err := viper.ReadConfig(bytes.NewBuffer(conf.AppJsonConfig)); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("no such config file")
		} else {
			// Config file was found but another error was produced
			log.Println("read config error")
		}
		log.Fatal(err) // 读取配置文件失败致命错误
	}
	models.InitDB()
}
func main() {
	//gin.SetMode(gin.ReleaseMode)
	port := viper.GetString("app.port")
	hub := ws.NewHub()
	go hub.Run()
	engine := routes.InitRoutes()
	if err := engine.Run(port); err != nil {
		println(err.Error())
	}
}
