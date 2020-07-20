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