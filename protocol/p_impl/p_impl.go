package p_impl

import (
	"chat_v3/protocol"
	//"chat_v3/protocol/creater"
	"chat_v3/protocol/handler"
	//"chat_v3/protocol/parser"
)

func init() {
	//protocol.ProtocParseMap = make(map[protocol.FrameType]protocol.ProtocParseFunc)
	protocol.ProtocHandleMap = make(map[protocol.FrameType]protocol.ProtocHandleFunc)
	//protocol.ProtocCreateMap = make(map[protocol.FrameType]protocol.ProtocCreateFunc)

	// 注册 handler
	protocol.RegisterNewHanleFunc(protocol.C_LOGIN, handler.HandleLogin)
	protocol.RegisterNewHanleFunc(protocol.C_JOIN, handler.HandleJoin)
	protocol.RegisterNewHanleFunc(protocol.C_MSG, handler.HandleRoomMsg)

	// 注册 creater
	//protocol.RegisterNewCreateFunc(protocol.C_LOGIN, creater.CreateLoginProto)
	//protocol.RegisterNewCreateFunc(protocol.C_JOIN, creater.CreateJoinProto)
	//protocol.RegisterNewCreateFunc(protocol.C_MSG, creater.CreateMsgProtoC)
	//protocol.RegisterNewCreateFunc(protocol.S_MSG, creater.CreateMsgProtoS)
	//protocol.RegisterNewCreateFunc(protocol.S_LOGIN, creater.CreateLoginOutcome)

	// 注册 parser
	//protocol.RegisterNewParseFunc(protocol.C_LOGIN, parser.ParseLoginProto)
	//protocol.RegisterNewParseFunc(protocol.C_JOIN, parser.ParseJoinProto)
	//protocol.RegisterNewParseFunc(protocol.C_MSG, parser.ParseMsgProto)
	//protocol.RegisterNewParseFunc(protocol.S_LOGIN, parser.ParseLoginOutcomeProto)
}
