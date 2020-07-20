package handler

import (
	"chat_v3/protocol"
	"chat_v3/client/ci"
)

func HandleLogin(p protocol.Protocol, c ci.Client) {
	c.HandleLogin(p)
}
