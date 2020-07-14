package client

import "net"

// 每Accept一个连接就创建一个Client，传入conn
// 之后进行登录等操作
// 登录成功后将该Client加入ClientManager中

type Client struct {
	id string
	clientChan chan string
    conn chan net.Conn
	online bool   // 已经登录为true
}


// 相当于构造函数
func NewClient(id string, conn net.Conn){
	c.id = id
	c.clientChan = make(chan string)
	c.conn
}