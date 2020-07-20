package main

import (
	"chat_v3/client"
	_ "chat_v3/protocol/p_impl"
	"fmt"
	"log"
	"net"
)

func main(){

	cm := client.NewClientManager()

	go cm.Broadcaster()
	go cm.Manager()

	listen, err := net.Listen("tcp", "localhost: 8000")
	if err != nil {
		log.Fatal(err)
		return
	}

	for{
		conn, err:= listen.Accept()
		fmt.Println("New Connection From ", conn.RemoteAddr().String())
		if err != nil {
			log.Println(err)
			continue
		}
		go cm.RunNewConnection(conn)
	}
}
