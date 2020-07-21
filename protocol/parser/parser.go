// 2020.7.21 将消息统一更换为NetMsg之后不再需要针对每一类型进行parse
// Parse函数采用 protocol.Parse 一个函数即可

package parser

import (
	"chat_v3/protocol"
	"encoding/json"
)

func ParseMsgProto(data []byte) protocol.Protocol {
	return string(data)
}

func ParseJoinProto(data []byte) protocol.Protocol{
	return string(data)
}

func ParseLoginProto(data []byte) protocol.Protocol{
	var info protocol.LogInfo
	_ = json.Unmarshal(data, &info)
	return info
}

func ParseLoginOutcomeProto(data []byte) protocol.Protocol{
	var success bool
	_ = json.Unmarshal(data, &success)
	return success
}