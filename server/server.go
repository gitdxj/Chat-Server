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
			
			// 最后一个人退出，关闭聊天室
			if len(clients) == 0{
				close(roomChannels[room].messages)
				close(roomChannels[room].leaving)
				close(roomChannels[room].entering)
				delete(roomChannels, room)
				return
			}
		}
	}
}

func handleConn(conn net.Conn){
	id := checkLogin(conn)
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
		} else if t == CMD_QUERY{
			rooms := getRooms()
			conn.Write([]byte(rooms))
		} else if t == CMD_JOIN{
			room = strings.Fields(str)[1]   // 获取房间名
			fmt.Println(id + "进入聊天室：" + room)
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

	// 从 客户端 接收消息 向同一聊天室内的人广播
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		roomChannels[room].messages <- string(buf[:n])
	}
	// 客户端终止进程，for循环结束
	// 通过leaving channel 通知broadcaster此人已离开
	roomChannels[room].leaving <- ch   // 离开房间



}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}


func checkLogin(conn net.Conn) (id string){
	buf := make([]byte, 1024)
	for {
		n, _ := conn.Read(buf)
		str := string(buf[:n])
		t, _ := parse(str)
		if t == CMD_LOGIN{   // 如果是登录
			context := strings.Fields(str)
			id := context[1]
			pswd := context[2]
			if checkIDPSWD(id, pswd){  // 登录不成功
				conn.Write([]byte("\\ok"))
				return id
			} else {
				conn.Write([]byte("\\wrong"))
				continue
			}
		} else {
			continue
		}
	}
}

func checkIDPSWD(id, pswd string) (ok bool) {
	if id == "dxj" && pswd == "123" || 
		id == "abc" && pswd == "123" {
		return true
	} else {
		return false
	}
}

func getRooms() (rooms string){
	if len(roomChannels) == 0 {
		return "未搜索到聊天室，创建一个吧\n"
	}
	for room, _ := range roomChannels{
		rooms += room + " "
	}
	return rooms+"\n"
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

