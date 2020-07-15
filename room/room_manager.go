package room
import (
	"chat/client"
	"sync"
)

type RoomManager struct {
	rooms map[string] *Room
	rwm sync.RWMutex
}

func NewRoomManager() *RoomManager{
	var rm RoomManager
	rm.rooms = make(map[string] *Room)
	return &rm
}

func (rm *RoomManager)getRoom(roomid string) (r *Room, ok bool){
	rm.rwm.RLock()
	r, ok = rm.rooms[roomid]
	rm.rwm.RUnlock()
	return r, ok
}

func (rm *RoomManager)HasRoom(roomid string) bool{
	_, ok := rm.getRoom(roomid)
	return ok
}

func (rm *RoomManager)CreateNewRoom(roomid string) {
	newRoom := NewRoom()

	rm.rwm.Lock()
	rm.rooms[roomid] = newRoom  // 添加房间
	rm.rwm.Unlock()
}

func (rm *RoomManager)RunRoom(roomid string, cm *client.ClientManager){
	rm.rwm.RLock()
	go rm.rooms[roomid].run(cm)
	rm.rwm.RUnlock()
}

func (rm *RoomManager)AddClientToRoom(id, roomid string){
	rm.rwm.RLock()
	go rm.rooms[roomid].addClient(id)
	rm.rwm.RUnlock()
}

func (rm *RoomManager)RemoveClientFromRoom(id, roomid string){
	rm.rwm.RLock()
	go rm.rooms[roomid].removeClient(id)
	rm.rwm.RUnlock()
}

func (rm *RoomManager)SendMsgToRoom(msg, roomid string){
	rm.rwm.RLock()
	go rm.rooms[roomid].sendMsgToRoom(msg)
	rm.rwm.RUnlock()
}

// GetAllRoomNames 返回现存的所有房间名
func (rm *RoomManager)GetAllRoomNames() (roomNames string) {
	if len(rm.rooms) == 0 {
		return "还没有任何房间，创建一个吧"
	}
	for r, _ := range rm.rooms {
		roomNames += r + " "
	}
	return roomNames
}
