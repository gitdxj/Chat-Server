package handler

import (
	"chat_v3/protocol"
	"chat_v3/client/ci"
)

func HandleJoin(p protocol.Protocol, c ci.Client){
	c.HandleJoin(p)
}
