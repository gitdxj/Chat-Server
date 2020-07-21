package main

import (
	"chat_v3/client"
	"chat_v3/config"
	_ "chat_v3/protocol/p_impl"
	"flag"
	"fmt"
	"log"
	"net"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.json", "服务器地址配置文件")
}

func main() {

	flag.Parse()

	cm := client.NewClientManager()
	go cm.Broadcaster()
	go cm.Manager()

	config.InitJsonConfig(configFile)
	addr := config.GetAddr()

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		conn, err := listen.Accept()
		fmt.Println("New Connection From ", conn.RemoteAddr().String())
		if err != nil {
			log.Println(err)
			continue
		}
		go cm.RunNewConnection(conn)
	}
}
