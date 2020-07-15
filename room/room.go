package room
import (
	"chat/client"
)

type Room struct {
	clients map[string] bool
	enterChan chan string
	leaveChan chan string
	msgChan chan string
}

func NewRoom() *Room {
	var r Room
	r.clients = make(map[string] bool)
	r.enterChan = make(chan string)
	r.leaveChan = make(chan string)
	r.msgChan = make(chan string)
	return &r
}

func (r *Room) addClient(id string){
	r.enterChan <- id
}
func (r *Room) removeClient(id string){
	r.leaveChan <- id
}

func (r *Room) sendMsgToRoom(msg string){
	r.msgChan <- msg
}

func (r *Room) closeRoom() {
	close(r.leaveChan)
	close(r.enterChan)
	close(r.msgChan)
}

func (r *Room) run(cm *client.ClientManager){
	for {
		select {
		case msg := <- r.msgChan:  // 从一个client接收到的消息
		// 转发给 聊天室内的所有用户
		for cli := range r.clients {
			go cm.SendMsgToClient(cli, msg)
		}

		case cli := <- r.enterChan:
			r.clients[cli] = true
		case cli := <- r.leaveChan:
			delete(r.clients, cli)
		}
	}
}
