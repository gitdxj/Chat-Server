package main

import (
	"chat/client"
	"chat/room"
	"log"
	"net"
)

func main(){

	rm := room.NewRoomManager()
	cm := client.NewClientManager()

	listen, err := net.Listen("tcp", "localhost: 8000")
	if err != nil {
		log.Fatal(err)
		return
	}

	for{
		conn, err:= listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go cm.RunNewConnection(conn, rm)
	}
}
