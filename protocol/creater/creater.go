package creater

import (
	"chat_v3/appsocket/byter"
	"chat_v3/protocol"
	"encoding/json"
)


func CreateMsgProtoC(p protocol.Protocol)  protocol.Protocol{
	msg := p.(string)
	return byter.CreateBS(protocol.C_MSG, []byte(msg))
}

func CreateMsgProtoS(p protocol.Protocol)  protocol.Protocol{
	msg := p.(string)
	return byter.CreateBS(protocol.S_MSG, []byte(msg))
}

func CreateJoinProto(p protocol.Protocol)  protocol.Protocol{
	roomid := p.(string)
	return byter.CreateBS(protocol.C_JOIN, []byte(roomid))
}

func CreateLoginProto(p protocol.Protocol)  protocol.Protocol{
	info := p.(protocol.LogInfo)
	infodata, _ := json.Marshal(info)
	return byter.CreateBS(protocol.C_LOGIN, infodata)
}

func CreateLoginOutcome( p protocol.Protocol)  protocol.Protocol{
	success := p.(bool)
	sucdata, _ := json.Marshal(success)
	return byter.CreateBS(protocol.S_LOGIN, sucdata)
}