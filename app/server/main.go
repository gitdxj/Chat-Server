package main

import (
	"chat_v3/client"
	"fmt"
	"log"
	"net"
)

func main(){

	cm := client.NewClientManager()

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
