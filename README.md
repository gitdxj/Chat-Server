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
1. 用户名密码发送到server端进行验证
2. 向服务器发送到信息不再简单使用string，要求TLV
3. 实现查询聊天室
4. 在某一聊天室中所有成员离开后将其资源释放

# 13.Jul
1. **√** 用户名密码发送到server端进行验证——在server端没有连接数据库，写死。账号有两个：dxj 123和 abc 123
```go
func checkIDPSWD(id, pswd string) (ok bool) {
	if id == "dxj" && pswd == "123" || 
		id == "abc" && pswd == "123" {
		return true
	} else {
		return false
	}
}
```
3. **√** 实现查询聊天室
4. **√** 在某一聊天室中所有成员离开后将其资源释放——关闭channels，相应map中把room删除

## todo 
1. 小 bug：发送消息后，发送者的消息他自己会收一遍，要不要改一下这里？
2. TLV搞一下

# 13.Jul review 修改意见
* 模块化——client, room, clientManager 和 roomManager 让代码结构清晰
* TLV机制—— 当我们使用net.Conn.Read读取字节流时，比如登录信息：\\login id pswd 有可能pswd还没有传输到，所以需要指定L(length)信息确保我们读取了一整条应用层的消息

# TLV 设计
* **T 4 Byte 类型**
* **L 4 Byte 长度**  protoBuf  --google varint
* **V L Byte 值**
## TLV解析算法
1. 读取 Tag（或Type）并使用 ntohl 将其转成主机字节序，指针偏移4；
2. 读取 Length ntohl** 将其转成主机字节序，指针偏移4；
3. 根据得到的长度读取 Value，若为 int、char、short、long 类型，将其转为主机字节序，指针偏移；若值为字符串，读取后指针偏移 Length；
4. 重复上述三步，继续读取后面的 TLV 单元。
```go
const TAG_SPECIFIER byte = 
```

# 16.Jul
* 重新设计结构——取消了Broadcaster，使用一个ClientManager来管理全部的连接
* TLV：封装conn.Read和conn.Write实现了appsocket/AppSocket，保证每次读取出来一个TLV结构
## todo
1. 用户名密码传输的时候使用json
2. 用户名和密码的验证尝试使用数据库
