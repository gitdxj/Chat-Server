package main

// 1. 所有人离开后房间仍然存在
// 2. 密码登录和

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type CommandType int8
const(
	CMD_QUERY CommandType = 0   // 查询聊天室
	CMD_JOIN CommandType = 1    // 加入聊天室
	CMD_LOGOUT CommandType = 2  // 注销
	CMD_LEAVE CommandType = 3   // 退出聊天室
	CMD_MSG CommandType = 4     // 纯消息
	CMD_HELP CommandType = 5
	CMD_LOGIN CommandType = 6
	CMD_EMPTY CommandType = 7
)


type client chan<-string
type channels struct{
	entering chan client
	leaving chan client
	messages chan string
}

// 每新创建一个聊天室，就在roomChannels中增加相应的房间，并初始化其中的三个channel
var roomChannels = make(map[string] channels)


func main(){
	listener, err := net.Listen("tcp", "localhost: 8000")
	if err != nil {
		log.Fatal(err)
	}

	for{
		conn, err := listener.Accept()
		if err != nil{
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}


func broadcaster(room string) {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-roomChannels[room].messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}

		case cli := <-roomChannels[room].entering:  // cli是一个 chan string 可以 赋给 type client chan<- string
			clients[cli] = true

		case cli := <-roomChannels[room].leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn){
	id, ok := checkLogin(conn)
	if ok {
		conn.Write([]byte("\\ok"))
	}
	fmt.Println(id, "登录成功")

	var room string
	ch := make(chan string)

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			//log.Println(err)
			return
		}
		str := string(buf[:n])
		t, _ := parse(str)
		if t == CMD_QUERY {
			conn.Write([]byte(getRooms()))
		} else if t == CMD_JOIN{
			room = strings.Fields(str)[1]   // 获取房间名
			fmt.Println(id + "有人进入房间了")
			// 每新创建一个聊天室，就在roomChannels中增加相应的，并初始化
			if _, exist := roomChannels[room]; exist == false{  // 没有该聊天室，创建新的
				var newRoom channels
				newRoom.entering = make(chan client)
				newRoom.leaving = make(chan client)
				newRoom.messages = make(chan string)
				roomChannels[room] = newRoom

				go broadcaster(room)  // 创建聊天室

				break
			} else {  // 已经有该聊天室了
				break
			}
		}
	}

	roomChannels[room].entering <- ch   // ch 通过 channel 通知 broadcaster

	go clientWriter(conn, ch)

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		roomChannels[room].messages <- string(buf[:n])
	}

	roomChannels[room].leaving <- ch   // 离开房间



}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}


func checkLogin(conn net.Conn) (id string, ok bool){
	buf := make([]byte, 1024)
	for {
		n, _ := conn.Read(buf)
		str := string(buf[:n])
		t, _ := parse(str)
		if t == CMD_LOGIN{   // 如果是登录
			context := strings.Fields(str)
			id := context[1]
			pswd := context[2]
			if checkIDPSWD(id, pswd){
				ok = true
				return id, ok
			}
		} else {
			continue
		}
	}
	ok = false
	return id, ok
}

func checkIDPSWD(id, pswd string) (ok bool) {
	return true
}

func getRooms() (rooms string){
	rooms = "abc 123"
	return rooms
}

func parse(str string) (t CommandType, bs []byte){
	context := strings.Fields(str)
	if len(context) == 0 {
		t = CMD_EMPTY
		return t, bs
	}
	t = CMD_MSG  // 默认为消息
	bs = []byte(str)
	if context[0] == "\\query"{
		t = CMD_QUERY
	}else if context[0] == "\\join"{
		t = CMD_JOIN
	}else if context[0] == "\\logout"{
		t = CMD_LOGOUT
	}else if context[0] == "\\leave"{
		t = CMD_LEAVE
	}else if context[0] == "\\help"{
		t = CMD_HELP
	}else if context[0] == "\\login"{
		t = CMD_LOGIN
	}
	return t, bs
}

