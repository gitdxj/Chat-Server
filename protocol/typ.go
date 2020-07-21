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

// 2020.7.21  使用统一的消息格式进行传递
type NetMsg struct {
	Id           string
	Pswd         string
	Msg          string
	RoomId       string
	LoginSuccess bool
}

