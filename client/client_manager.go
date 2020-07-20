package client

import (
	"chat_v3/protocol"
	"fmt"
	"log"
	"net"
	"sync"
)

type CMChans struct {
	addClientChan    chan *Client
	removeClientChan chan *Client
	broadcastChan    chan protocol.BroadcastMsg
	joinChan         chan protocol.JoinMsg
}

type ClientManager struct {
	clients   map[string]*Client
	rooms     map[string]map[string]bool // rooms[room][user] = true
	clientRWM sync.RWMutex
	roomRWM   sync.RWMutex
	cmChans   CMChans
}

func NewClientManager() *ClientManager {
	var cm ClientManager
	cm.clients = make(map[string]*Client)
	cm.rooms = make(map[string]map[string]bool)
	cm.cmChans.addClientChan = make(chan *Client)
	cm.cmChans.removeClientChan = make(chan *Client)
	cm.cmChans.broadcastChan = make(chan protocol.BroadcastMsg)
	cm.cmChans.joinChan = make(chan protocol.JoinMsg)
	return &cm
}

// AddClient 向clients中添加新的Client
func (cm *ClientManager) AddClient(c *Client) bool {
	if !c.online {
		return false
	}
	cm.clients[c.id] = c
	return true
}

// HasRoom 判断cm中是否存在roomid
func (cm *ClientManager) HasRoom(roomid string) bool {
	_, has := cm.rooms[roomid]
	return has
}

// CreateRoom 创建名为roomid的房间
func (cm *ClientManager) CreateRoom(roomid string) {
	cm.rooms[roomid] = make(map[string]bool)
}

// RemoveClient 将用户名=id的用户从clients中删除
func (cm *ClientManager) RemoveClient(id string) {
	delete(cm.clients, id)
}

// AddClientToRoom 将用户添加到房间
func (cm *ClientManager) AddClientToRoom(id, roomid string) {
	if !cm.HasRoom(roomid) {
		cm.CreateRoom(roomid)
	}
	cm.rooms[roomid][id] = true
}

// RemoveClientFromRoom 将用户从房间中删除
func (cm *ClientManager) RemoveClientFromRoom(id, roomid string) {
	delete(cm.rooms[roomid], id)
	fmt.Println("已将", id, " 从房间", roomid, "删除")
}

// clientChan用来接收从Room广播的消息，Room通过调用这个函数来将msg传入clientChan
func (cm *ClientManager) SendMsgToClient(id, msg string) {
	cli, ok := cm.clients[id]
	if !ok {
		log.Println("向不存在的用户", cli, "发送了消息")
		return
	}
	cli.sendMsg(msg)
}

// BroadcastInRoom 向房间内的所有用户广播一条消息
func (cm *ClientManager) BroadcastInRoom(roomid, msg string) {
	clientIDs := cm.rooms[roomid]
	for id, _ := range clientIDs {
		fmt.Println("向", id, "发送了消息", msg)
		cm.SendMsgToClient(id, msg)
	}
}

// GetAllRoomNames 返回现在所有的房间名称
func (cm *ClientManager) GetAllRoomNames() (str string) {
	for roomid, _ := range cm.rooms {
		str += roomid + " "
	}
	if str == "" {
		return "现在还没有房间，创建一个吧"
	}
	return str
}

func (cm *ClientManager) Broadcaster() {
	for bcMsg := range cm.cmChans.broadcastChan {
		log.Println("向房间", bcMsg.RoomId, "转发： ", bcMsg.Msg)
		go cm.BroadcastInRoom(bcMsg.RoomId, bcMsg.Msg)
	}
}

func (cm *ClientManager) Manager() {
	for {
		select {
		case c := <-cm.cmChans.addClientChan:
			log.Println("AddClient")
			cm.AddClient(c)
		case c := <-cm.cmChans.removeClientChan:
			log.Println("RemoveClientFromRoom")
			cm.RemoveClientFromRoom(c.GetId(), c.GetRoomId())
		case jm := <-cm.cmChans.joinChan:
			log.Println("AddClientToRoom")
			cm.AddClientToRoom(jm.ID, jm.RoomId)
		}
	}
}

func (cm *ClientManager) RunNewConnection(conn net.Conn){
	c := NewClient(conn, cm)
	syn := make(chan struct{})
	go c.Run(syn)
	<- syn
	fmt.Println("RunNewConnection Ends")

	if c.GetRoomId() != "" {
		cm.RemoveClientFromRoom(c.GetId(), c.GetRoomId())
	}
	if c.GetId() != "" {
		cm.RemoveClient(c.GetId())
	}
}

