package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

/*
	Status 是client端用户所处状态
	不同状态下用户的输入采取的处理不同
	在命令模式下有查询、加入聊天室等一系列指令
	而输入其他文字会直接被丢弃
	在聊天模式下，输入的文字会被发送到服务器
	同时也有退出等指令
*/
type Status int8
const(
	STAT_LOGIN Status = 0
	STAT_CMD Status = 1
	STAT_IN_ROOM Status = 2
)
var currentStatus Status = STAT_LOGIN


type CommandType int8
const(
	CMD_QUERY CommandType = 0   // 查询聊天室
	CMD_JOIN CommandType = 1    // 加入聊天室
	CMD_LOGOUT CommandType = 2  // 注销
	CMD_LEAVE CommandType = 3   // 退出聊天室
	CMD_MSG CommandType = 4     // 纯消息
	CMD_HELP CommandType = 5    // 帮助，显示命令提示
	CMD_LOGIN CommandType = 6   // 登录
	CMD_EMPTY CommandType = 7   // 空命令，为了解决在parse函数中传入了一个空字符串的问题
)

const allCommand string= "\\query roomid  查找指定房间\n" +
"\\join roomid 加入房间，若无此房间则删除\n" +
"\\logout 注销并退出\n" +
"\\help 显示所有命令\n" +
"在聊天室内：\n" +
"\\leave 退出聊天室\n" +
"\\logout 注销并退出\n"

func main(){
	conn, err := net.Dial("tcp", "localhost: 8000")
	if err != nil{
		log.Println(err)
	}
	id := login(conn)  // 通过登录
	currentStatus = STAT_CMD

	showAllCommand()  // help - 显示提示信息

	// done := make(chan struct{})  // 同步发送消息例程和接收消息（主例程）

	// 发送消息例程：
	// 用户在输入框中输入，可能是聊天框命令，也可能是普通消息
	// 输入内容，经过处理后，得到字节流
	go func(){
		defer conn.Close()
		input := bufio.NewScanner(os.Stdin)
		input.Split(bufio.ScanLines)
		for input.Scan() {
			var str string = input.Text()   // 获取用户命令行输入
			cmdtype, bs := parse(str)
			if currentStatus == STAT_CMD{  // 如果现在正命令模式
				if cmdtype == CMD_HELP {
					showAllCommand()
					continue
				} else if cmdtype == CMD_QUERY {
					conn.Write(bs)
					continue
				} else if cmdtype == CMD_JOIN{
					currentStatus = STAT_IN_ROOM
					conn.Write(bs)
					continue
				} else if cmdtype == CMD_LOGOUT{
					return
				} else {
					continue
				}

			} else if currentStatus == STAT_IN_ROOM{  // 如果现在正在聊天室内
				if cmdtype == CMD_LOGOUT{
					return
				} else if cmdtype == CMD_LEAVE{
					conn.Write(bs)
					continue
				} else if cmdtype == CMD_MSG{
					bs = []byte(id+":  "+string(bs))
					conn.Write(bs)
					continue
				} else {
					continue
				}
			}
		}
	}()

	// 接收消息，接收的消息一律打印即可
	_, err = io.Copy(os.Stdout, conn)   // Copy会一直阻塞到产生错误或EOF
	if err != nil {
		// log.Println(err)   // 当我们在发送消息例程中把conn关闭的时候，这里一定会报错
		return
	}

}

// login 登录
func login(conn net.Conn) (id string){
	for {
		id, pswd := idInput()     // 键入用户名和密码
		if ok := checkID(id, pswd, conn); ok {  // 发送到服务器进行验证
			return id                           // 验证成功，返回用户id
		} else {                                // 用户名密码不正确重新输入
			fmt.Println("输入的用户名和密码不正确")
			continue
		}
		//fmt.Println(id, pswd)
	}
}


// idInpit 键盘输入用户名和密码
func idInput() (id, pswd string){
	fmt.Print("请输入用户名：")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	id = input.Text()

	fmt.Print("请输入密码：")
	input.Scan()
	pswd = input.Text()

	return id, pswd
}


func parse(str string) (t CommandType, bs []byte){
	context := strings.Fields(str)
	if len(context) == 0 {
		t = CMD_EMPTY
		return t, bs
	}
	t = CMD_MSG  // 默认为消息
	bs = []byte(str)  // 在这里为了简单就直接转为字节流了 没有TLV
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
	} else if context[0] == "\\login"{
		t = CMD_LOGIN
	}
	return t, bs
}

// checkID 查询用户名和密码是否正确
func checkID(id, pswd string, conn net.Conn) bool{
	bs := []byte("\\login " + id + " " + pswd + "\n")
	conn.Write(bs)
	check := make([]byte, 1024)
	conn.Read(check)
	//context := strings.Fields(string(check))
	if "\\ok" == string(check[:3]){
		return true
	} else {
		return false
	}
}

// 查询聊天室 之后看看能不能加入通配查询的功能
func queryRoom(roomID string, conn net.Conn) (roomList []string) {
	rooms := []string{"abc", "zulong", "123"}
	return rooms
}

func showAllCommand(){
	fmt.Print(allCommand)
}