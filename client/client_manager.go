package client

import (
	"chat_v3/protocol"
	"fmt"
	"log"
	"net"
	"sync"
)

// cm作为一个全局变量
// 每Accept一个连接，go cm.Login()


type ClientManager struct {
	clients map[string] *Client
	rooms map[string]  map[string] bool // rooms[room][user] = true
	clientRWM sync.RWMutex
	roomRWM sync.RWMutex
}


func NewClientManager() *ClientManager{
	var cm ClientManager
	cm.clients = make(map[string] *Client)
	cm.rooms = make(map[string] map[string] bool)
	return &cm
}

// AddClient 向clients中添加新的Client
func (cm *ClientManager) AddClient(c *Client) bool{
	if !c.online {
		return false
	}
	cm.clients[c.id] = c
	return true
}

// HasRoom 判断cm中是否存在roomid
func (cm *ClientManager) HasRoom(roomid string) bool{
	_, has := cm.rooms[roomid]
	return has
}

// CreateRoom 创建名为roomid的房间
func (cm *ClientManager) CreateRoom(roomid string) {
	cm.rooms[roomid] = make(map[string] bool)
}

// RemoveClient 将用户名=id的用户从clients中删除
func (cm *ClientManager) RemoveClient(id string) {
	delete(cm.clients, id)
}

// AddClientToRoom 将用户添加到房间
func (cm *ClientManager) AddClientToRoom(id, roomid string){
	if !cm.HasRoom(roomid) {
		cm.CreateRoom(roomid)
	}
	cm.rooms[roomid][id] = true
}

// RemoveClientFromRoom 将用户从房间中删除
func (cm *ClientManager) RemoveClientFromRoom(id, roomid string){
	delete(cm.rooms[roomid], id)
	fmt.Println("已将", id, " 从房间", roomid, "删除")
}

// clientChan用来接收从Room广播的消息，Room通过调用这个函数来将msg传入clientChan
func (cm *ClientManager) SendMsgToClient(id, msg string){
	cli, ok := cm.clients[id]
	if !ok {
		log.Println("向不存在的用户", cli, "发送了消息")
		return
	}
	cli.sendMsg(msg)
}

// BroadcastInRoom 向房间内的所有用户广播一条消息
func (cm *ClientManager) BroadcastInRoom(roomid, msg string){
	clientIDs := cm.rooms[roomid]
	for id, _ := range clientIDs{
		fmt.Println("向", id, "发送了消息", msg)
		cm.SendMsgToClient(id, msg)
	}
}

// GetAllRoomNames 返回现在所有的房间名称
func (cm *ClientManager) GetAllRoomNames() (str string){
	for roomid, _ := range cm.rooms {
		str += roomid + " "
	}
	if str == "" {
		return "现在还没有房间，创建一个吧"
	}
	return str
}

// 每Accept一个连接，就go cm.RunNewConnection()
func (cm *ClientManager) RunNewConnection(conn net.Conn) {
	c := NewClient(conn)
	for {
		ft, val, err := c.Read()
		fmt.Println("New MSG TYPE = ", uint32(ft), string(val))
		if err != nil {
			return
		}
		if !c.online { // 若还没登录
			if ft == protocol.T_LOGIN {
				id, pswd := protocol.ParseLogInfo(val)
				if CheckLogIn(id, pswd) {  // 登录成功
					c.id = id
					c.online = true
					cm.AddClient(c)
					c.SendLogInSuccess()
					fmt.Println(c.id, "登录成功")
					defer cm.RemoveClient(c.id)
				} else {
					c.SendLoginFail()
				}
			}
		} else {                       // 若已经登录
			if ft == protocol.T_JOIN { // 类型为加入房间
				roomid := protocol.ParseRoomId(val)
				fmt.Println("加入房间:", roomid)
				cm.AddClientToRoom(c.id, roomid)
				c.roomid = roomid
				go c.recvFromRoomAndSend()   // 从clientChan 收取同一房间的广播消息然后发给客户端
				defer cm.RemoveClientFromRoom(c.id, c.roomid)
			} else if ft == protocol.T_QUERY { // 类型为房间查询
				roomNames := cm.GetAllRoomNames()
				_, err = c.as.WriteAppFrame(protocol.CreateMsgBS(roomNames))
				if err != nil {
					log.Println(err)
				}
			} else if ft == protocol.T_MSG { // 类型为消息
				// 如果c在某一房间内，就向这一房间的用户广播
				if c.roomid != "" {
					fmt.Println("所在的房间为 ", c.roomid)
					cm.BroadcastInRoom(c.roomid, string(val))
				}
			}
		}
	}
}

