# Chat Server
# 服务器端设计

## goroutine
* 每一个接入服务器的用户——go handleConn  实现多用户
* 每一个聊天室——go broadcaster  用来向同一聊天室的用户广播消息

## Channel
```go
 // 每个用户接入时 hanleConn中 创建一个client，用来接收 broadcaster 广播的消息
type client chan<-string 
type channels struct{
    // 用户进入某聊天室时，通过 entering 向 broadcaster 传入 client channel
    entering chan client
    // 用户离开某聊天室时，通过 leaving 向 broadcaster 传入 client channel
    leaving chan client
    // 用户发送的消息 通过 messages 传入 broadcaster
	messages chan string
}
```
在 broadcaster 内三个 channel 使用 select 并列

# 11.Jul
## todo
* 用户名密码发送到server端进行验证
* 向服务器发送到信息不再简单使用string，要求TLV
* 实现查询聊天室
* 在某一聊天室中所有成员离开后将其资源释放