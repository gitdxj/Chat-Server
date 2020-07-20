package protocol


type BroadcastMsg struct {
	Msg    string
	RoomId string
}

type JoinMsg struct {
	ID     string
	RoomId string
}

type LogInfo struct {
	Id   string
	Pswd string
}

