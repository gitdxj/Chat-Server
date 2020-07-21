package protocol

import (
	"chat_v3/client/ci"
	"encoding/json"
	"fmt"
)

type FrameType uint32

const (
	S_LOGIN FrameType = 101 // 登录
	C_LOGIN FrameType = 102
	S_JOIN  FrameType = 201 // 加入聊天室
	C_JOIN  FrameType = 202
	S_QUERY FrameType = 301 // 查询聊天室
	C_QUERY FrameType = 302
	S_MSG   FrameType = 401 // 消息
	C_MSG   FrameType = 402 // 登录失败
)

type Protocol interface {
}

//type ProtocParseFunc func(data []byte) Protocol
type ProtocHandleFunc func(p Protocol, c ci.Client)
type ProtocCreateFunc func(p Protocol) Protocol

//var ProtocParseMap map[FrameType]ProtocParseFunc
var ProtocHandleMap map[FrameType]ProtocHandleFunc
var ProtocCreateMap map[FrameType]ProtocCreateFunc



//func RegisterNewParseFunc(ft FrameType, f ProtocParseFunc) {
//	ProtocParseMap[ft] = f
//}

func RegisterNewHanleFunc(ft FrameType, f ProtocHandleFunc) {
	ProtocHandleMap[ft] = f
}

func RegisterNewCreateFunc(ft FrameType, f ProtocCreateFunc) {
	ProtocCreateMap[ft] = f
}

//func Parse(ft FrameType, data []byte) Protocol{
//	return ProtocParseMap[ft](data)
//}

// 所有的消息都用统一的结构后直接Unmarshal就可以了
func Parse(data []byte) Protocol{
	var nm NetMsg
	err := json.Unmarshal(data, &nm)
	if err != nil {
		fmt.Println("Parse:", err)
	}
	return nm
}

func Handle(ft FrameType, p Protocol, c ci.Client){
	ProtocHandleMap[ft](p, c)
}

//func Create(ft FrameType, p Protocol) Protocol{
//	return ProtocCreateMap[ft](p)
//}

func Create(ft FrameType, p Protocol) Protocol{
	netMsg := p.(NetMsg)
	val, _ := json.Marshal(netMsg)
	return CreateBS(ft, val)
}



//type NewProtocolFunc func() interface{}

//var protoCreatMap map[FrameType]NewProtocolFunc
//func registerNewFunc(ft FrameType, f NewProtocolFunc) {
//	protoCreatMap[ft] = f
//}

//func ParseProtoc(typ FrameType, buf []byte) interface{} {
//	f, ok := protoCreatMap[typ]
//	p := f()
//	json.Unmarshal(buf, p)
//	handler, ok := handmap[typ]
//	handler(p)
//	return p
//}
//func handlerLogin(p interface{}) error {
//	p, ok := p.(*LogInfo)
//}

