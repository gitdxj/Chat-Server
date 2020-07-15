package client

import (
	"chat/appsocket"
	"log"
	"net"
)

type Client struct {
	id string
	clientChan chan string
    as appsocket.AppSocket
	online bool   // 已经登录为true
	roomid string
}

// 相当于构造函数
func NewClient(conn net.Conn) * Client{
	var c Client
	c.clientChan = make(chan string)
	c.as = *appsocket.NewAppSocket(conn)
	c.online = false
	return &c
}

func (c *Client)SendMsgToClient(str string){
	c.clientChan <- str
}

func (c *Client)Read() (appsocket.FrameType, []byte, error){
	return c.as.ReadAppFrame()
}

func (c *Client)Write(bs []byte) (n int, err error){
	return c.as.WriteAppFrame(bs)
}

func (c *Client)SendLogInSuccess() {
	_, err := c.Write(appsocket.CreateBS(appsocket.T_LOGIN_SUCCESS, ""))
	if err != nil {
		log.Println("SendLogInSuccess", err)
		return
	}
}

func (c *Client)SendLoginFail(){
	_, err := c.Write(appsocket.CreateBS(appsocket.T_LOGIN_FAIL, ""))
	if err != nil {
		log.Println("SendLogInSuccess", err)
		return
	}
}

func CheckLogIn(id, pswd string) bool{
	return true
}

// recvFromRoomAndSend 从聊天室接收消息并发送到客户端
func (c *Client)recvFromRoomAndSend(){
	for msg := range c.clientChan {
		_, err := c.as.WriteAppFrame(appsocket.CreateMsgBS(msg))
		if err != nil {
			log.Println("recvFromRoomAndSend", err)
			return
		}
	}
}

func (c *Client)getId() string {
	return c.id
}