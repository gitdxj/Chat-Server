package main

import (
	"chat_v3/config"
	"flag"
)

var configFile string

func init(){
	flag.StringVar(&configFile, "config", "config.json", "服务器地址配置文件")
}

func main(){

	flag.Parse()

	config.InitJsonConfig(configFile)
	serverAddr := config.GetAddr()

	Run(serverAddr)
}