package ci

type Client interface {
	//SetId(string)
	//SetOnline(bool)
	//SendLoginOutcome(bool)
	//
	//AddToCM()
	//BroadcastMsg(protocol.BroadcastMsg)
	//JoinRoom(protocol.JoinMsg)

	HandleLogin(i interface{})
	HandleBroadcastMsg(i interface{})
	HandleJoin(i interface{})
}
