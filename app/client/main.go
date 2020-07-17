package main

import (
	"bufio"
	"chat_v3/appsocket"
	"chat_v3/protocol"
	"fmt"
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
	as := appsocket.NewAppSocket(conn)
	id, err := login(as)  // 通过登录
	if err != nil {
		log.Println("main: login", err)
		return
	}


	currentStatus = STAT_CMD

	showAllCommand()  // help - 显示提示信息

	// done := make(chan struct{})  // 同步发送消息例程和接收消息（主例程）

	// 发送消息 goroutine:
	// 用户在输入框中输入，可能是聊天框命令，也可能是普通消息
	// 输入内容，经过处理后，得到字节流
	go func(){
		defer as.Close()
		input := bufio.NewScanner(os.Stdin)
		input.Split(bufio.ScanLines)
		for input.Scan() {
			var str string = input.Text()   // 获取用户命令行输入
			cmdtype, bs := parseInput(str, id)
			if currentStatus == STAT_CMD{  // 如果现在正命令模式
				if cmdtype == CMD_HELP {
					showAllCommand()
					continue
				} else if cmdtype == CMD_QUERY {
					as.WriteAppFrame(bs)
					continue
				} else if cmdtype == CMD_JOIN{
					currentStatus = STAT_IN_ROOM
					as.WriteAppFrame(bs)
					continue
				} else if cmdtype == CMD_LOGOUT{
					return
				} else {
					continue
				}

			} else if currentStatus == STAT_IN_ROOM{  // 如果现在正在聊天室内
				if cmdtype == CMD_LOGOUT{
					return
				} else if cmdtype == CMD_MSG{
					as.WriteAppFrame(bs)
					continue
				} else {
					continue
				}
			}
		}
	}()

	for {
		_, val, err := as.ReadAppFrame()
		if err != nil {
			break
		}
		fmt.Println(string(val))
	}

}

// login 登录
func login(as *appsocket.AppSocket) (id string, err error){
	for {
		id, pswd := idInput()     // 键入用户名和密码
		ok, err := checkID(id, pswd, as)
		fmt.Println("checkID ....")
		if err != nil {
			log.Println("login", err)
			return id, err
		}
		if ok {
			return id, nil
		} else {
			continue
		}
	}
}

// checkID 查询用户名和密码是否正确
func checkID(id, pswd string, as *appsocket.AppSocket) (ok bool, err error){
	bs := protocol.CreateLoginBS(id, pswd)
	_, err = as.WriteAppFrame(bs)
	if err != nil {
		log.Fatal("checkID", err)
	}
	fmt.Println("ID and Password sent to server")
	ft, _, err := as.ReadAppFrame()
	if err != nil {
		return false, err
	}

	if ft == protocol.T_LOGIN_SUCCESS {
		return true, nil
	} else if ft == protocol.T_LOGIN_FAIL {
		return false, nil
	}
	return ok, err
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


func parseInput(str, id string) (t CommandType, bs []byte){
	// 处理输入了条空行
	if str == "" {
		t = CMD_EMPTY
		return t, bs
	}

	context := strings.Fields(str)
	tag := context[0]
	switch tag {
	case "\\query":
		return CMD_QUERY,  protocol.CreateQueryBS()
	case "\\join":
		if len(context) < 2 {
			return CMD_EMPTY, bs
		}
		return CMD_JOIN, protocol.CreateJoinBS(context[1])
	case "\\logout":
		return CMD_LOGOUT, bs
	case "\\help":
		return CMD_HELP, bs
	default:
		str = id + ": " + str
		return CMD_MSG, protocol.CreateMsgBS(str)
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