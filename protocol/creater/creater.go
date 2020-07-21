// 2020.7.21 将消息统一更换为NetMsg之后不再需要针对每一类型进行序列化
// Create函数采用 protocol.Create 一个函数即可

package creater

import (
	"chat_v3/protocol"
	"encoding/json"
)


func CreateMsgProtoC(p protocol.Protocol)  protocol.Protocol{
	msg := p.(string)
	return protocol.CreateBS(protocol.C_MSG, []byte(msg))
}

func CreateMsgProtoS(p protocol.Protocol)  protocol.Protocol{
	msg := p.(string)
	return protocol.CreateBS(protocol.S_MSG, []byte(msg))
}

func CreateJoinProto(p protocol.Protocol)  protocol.Protocol{
	roomid := p.(string)
	return protocol.CreateBS(protocol.C_JOIN, []byte(roomid))
}

func CreateLoginProto(p protocol.Protocol)  protocol.Protocol{
	info := p.(protocol.LogInfo)
	infodata, _ := json.Marshal(info)
	return protocol.CreateBS(protocol.C_LOGIN, infodata)
}

func CreateLoginOutcome( p protocol.Protocol)  protocol.Protocol{
	success := p.(bool)
	sucdata, _ := json.Marshal(success)
	return protocol.CreateBS(protocol.S_LOGIN, sucdata)
}