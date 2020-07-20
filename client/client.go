package client

import (
	"chat_v3/appsocket"
	"chat_v3/protocol"
	"fmt"

	"encoding/json"
	"log"
	"net"
)

type Client struct {
	id         string
	clientChan chan string
	as         appsocket.AppSocket
	online     bool   // 已经登录为true
	roomid     string // 在哪个房间
	cmChans    *CMChans
}

// 相当于构造函数
func NewClient(conn net.Conn, cm *ClientManager) *Client {
	var c Client
	c.clientChan = make(chan string)
	c.as = *appsocket.NewAppSocket(conn)
	c.online = false
	c.cmChans = &cm.cmChans
	return &c
}

func (c *Client) sendMsg(str string) {
	c.clientChan <- str
}

func (c *Client) SetId(id string) {
	c.id = id
}

func (c *Client) SetRoom(roomid string) {
	c.roomid = roomid
}

func (c *Client) SetOnline(online bool) {
	c.online = online
}

func (c *Client) GetId() string {
	return c.id
}

func (c *Client) GetRoomId() string {
	return c.roomid
}

func (c *Client) Read() (protocol.FrameType, []byte, error) {
	return c.as.ReadAppFrame()
}

func (c *Client) Write(bs []byte) (n int, err error) {
	return c.as.WriteAppFrame(bs)
}

func (c *Client) SendLoginOutcome(success bool){
	c.Write(protocol.Create(protocol.S_LOGIN, success).([]byte))
}

func CheckLogIn(id, pswd string) bool {
	return true
}

func CheckLogInfo(info protocol.LogInfo) bool {
	return true
}

func parseLogin(buf []byte) interface{} {
	var p interface{}
	err := json.Unmarshal(buf, p)
	if err != nil {
		log.Fatal("parseLoglin", err)
	}
	return p
}

func parseMsg(buf []byte) interface{} {
	return string(buf)
}

func parseJoin(buf []byte) interface{} {
	return string(buf)
}

func (c *Client) AddToCM() {
	c.cmChans.addClientChan <- c
}

func (c *Client) BroadcastMsg(bm protocol.BroadcastMsg) {
	c.cmChans.broadcastChan <- bm
}

func (c *Client) JoinRoom(jm protocol.JoinMsg) {
	c.cmChans.joinChan <- jm
}

func (c *Client) HandleLogin(i interface{}) {
	info := i.(protocol.LogInfo)
	if CheckLogInfo(info) {
		c.SetId(info.Id)
		c.SetOnline(true)
		c.SendLoginOutcome(true)
		c.AddToCM()
	} else {
		c.SendLoginOutcome(false)
	}
}

func (c *Client) HandleBroadcastMsg(i interface{}) {
	msg := i.(string)
	bm := protocol.BroadcastMsg{
		msg,
		c.roomid,
	}
	c.BroadcastMsg(bm)
}

func (c *Client) HandleJoin(i interface{}) {
	roomid := i.(string)
	fmt.Println(c.id, "正在Handle Join", roomid)
	jm := protocol.JoinMsg{
		c.id,
		roomid,
	}
	c.SetRoom(roomid)
	go c.JoinRoom(jm)
	fmt.Println("用户", c.id, "加入了房间", roomid)
}

// recvFromRoomAndSend 从聊天室接收消息并发送到客户端
func (c *Client) recvFromRoomAndSend() {
	for msg := range c.clientChan {
		_, err := c.as.WriteAppFrame(protocol.Create(protocol.S_MSG, msg).([]byte))
		if err != nil {
			log.Println("recvFromRoomAndSend", err)
			return
		}
	}
}

func (c* Client) Run(syn chan struct{}){
	go c.recvFromRoomAndSend()
	for{
		ft, val, err := c.Read()
		fmt.Printf("%d   %s\n",ft , val)
		if err != nil {
			break
		}
		p := protocol.Parse(ft, val)
		protocol.Handle(ft, p, c)
	}
	syn <- struct{}{}
}

