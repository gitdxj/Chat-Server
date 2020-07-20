package handler

import (
	"chat_v3/client/ci"
	"chat_v3/protocol"
)

func HandleRoomMsg (p protocol.Protocol, c ci.Client){
	c.HandleBroadcastMsg(p)
}
