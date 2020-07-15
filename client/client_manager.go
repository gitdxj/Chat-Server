package client

import (
	"chat/appsocket"
	"chat/room"
	"fmt"
	"log"
	"net"
	"sync"
)

// cm作为一个全局变量
// 每Accept一个连接，go cm.Login()


type ClientManager struct {
	clients map[string] *Client
	rwm sync.RWMutex
}

func NewClientManager() *ClientManager{
	var cm ClientManager
	cm.clients = make(map[string] *Client)
	return &cm
}

func (cm *ClientManager) AddClient(c *Client) bool{
	if !c.online {
		return false
	}
	cm.rwm.Lock()
	cm.clients[c.id] = c
	cm.rwm.Unlock()
	return true
}

func (cm *ClientManager) RemoveClient(id string) {
	cm.rwm.Lock()
	delete(cm.clients, id)
	cm.rwm.Unlock()
}

// clientChan用来接收从Room广播的消息，Room通过调用这个函数来将msg传入clientChan
func (cm *ClientManager) SendMsgToClient(id, msg string){
	cli, ok := cm.clients[id]
	if !ok {
		log.Println("向不存在的用户发送了消息")
		return
	}
	cli.clientChan <- msg
}

// 每Accept一个连接，就go cm.RunNewConnection()
func (cm *ClientManager)RunNewConnection(conn net.Conn, rm *room.RoomManager) {
	c := NewClient(conn)
	for {
		ft, val, err := c.Read()
		if err != nil {
			break
		}
		if !c.online { // 若还没登录
			if ft == appsocket.T_LOGIN {
				id, pswd := appsocket.ParseLogInfo(val)
				if CheckLogIn(id, pswd) {  // 登录成功
					c.id = id
					c.online = true
					cm.AddClient(c)
					c.SendLogInSuccess()
					fmt.Println(c.id, "登录成功")
				} else {
					c.SendLoginFail()
				}
			}
		} else if c.online { // 若已经登录
			if ft == appsocket.T_JOIN {
				roomid := appsocket.ParseRoomId(val)
				if !rm.HasRoom(roomid) {
					rm.CreateNewRoom(roomid)
					rm.RunRoom(roomid, cm)
				}
				rm.AddClientToRoom(c.id, roomid)
				c.roomid = roomid
				defer func(){
					rm.RemoveClientFromRoom(c.id, roomid)
					cm.RemoveClient(c.id)
				}()
				go c.recvFromRoomAndSend()
			} else if ft == appsocket.T_QUERY {
				roomNames := rm.GetAllRoomNames()
				_, err = c.as.WriteAppFrame([]byte(roomNames))
				if err != nil {
					log.Println(err)
				}
			} else if ft == appsocket.T_MSG {
				rm.SendMsgToRoom(string(val), c.roomid)
			}
		}
	}
}
