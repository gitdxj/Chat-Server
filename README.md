# 聊天室
# 运行
1. 在`app/client/config.json`和`app/server/config.json`中设置`Dial`和`Listen`的地址和端口
2. `app/server`目录下`go run main.go`
3. `app/client`目录下`go run main.go app.go`
# 客户端
登录时使用任意用户名密码即可，服务器不进行检查  
**客户端指令：**
* /join ROOMID 加入名称为ROOMID的房间
* /logout
* /help 

`\query`查询房间的指令在17.Jul更改代码结构后就删除了


# 服务器端设计

## goroutine
* 每一个接入服务器的用户——`go handleConn`  实现多用户
* 每一个聊天室——`go broadcaster ` 用来向同一聊天室的用户广播消息

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
在 `broadcaster` 内三个 `channel` 使用 `select` 并列

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
* 模块化——`client`, `room`, `clientManager` 和 `roomManager` 让代码结构清晰
* TLV机制—— 当我们使用`net.Conn.Read`读取字节流时，比如登录信息：\\login id pswd 有可能pswd还没有传输到，所以需要指定L(length)信息确保我们读取了一整条应用层的消息

# TLV 设计
* **T 4 Byte 类型**
* **L 4 Byte 长度**  protoBuf  --google varint
* **V L Byte 值**
## TLV解析算法
1. 读取 `Tag`（或`Type`）并使用 `ntohl` 将其转成主机字节序，指针偏移4；
2. 读取 `Length` ntohl** 将其转成主机字节序，指针偏移4；
3. 根据得到的长度读取 `Value`，若为 `int`、`char`、`short`、`long` 类型，将其转为主机字节序，指针偏移；若值为字符串，读取后指针偏移 `Length`；
4. 重复上述三步，继续读取后面的 TLV 单元。
```go
const TAG_SPECIFIER byte = 
```

# 16.Jul
* 重新设计结构——取消了`Broadcaster`，使用一个`ClientManager`来管理全部的连接
* TLV：封装`conn.Read`和`conn.Write`实现了`appsocket/AppSocket`，保证每次读取出来一个TLV结构
## todo
1. 用户名密码传输的时候使用json
2. 用户名和密码的验证尝试使用数据库

# 17.Jul
## todo
* 重新修改代码结构
* **Protocol 重新定义**
* Protocol定义上把Server和Client分开
* **封装 封装 封装**

# 20.Jul

1. 之前的代码比较繁琐，现在使用一个Protocol接口进行了封装，一开始给我示例的时候我也不知道这里具体要怎么做，只用一个接口如何进行封装呢？ Protocol的定义非常简单，就是  
```go
var Protocal interface{}
```  
**那么一切皆可为Protocal，大师我悟了**

具体是使用类似于函数指针的方法，具体指令内容->`Create`->字节流->`Parse`->结构化数据->`Handle` 处理，其中`Create` `Parse` `Handle`通过类似于函数指针的方式，根据相应的`FrameType`来选择相应的处理函数。

## todo
* **统一的序列化和反序列化**：在这个版本中，每种消息都使用的不同的数据结构，比如加入聊天室`\join roomid`其中`roomid`是一个`string`序列化时直接被我转换成`[]byte`了，而登录`\login info`， `info`是一个`LogInfo`结构体
```go
type LogInfo struct {
	Id   string,
	Pswd string
}
```
对这种消息进行序列化和反序列化使用的是`json.Marshal`和`json.Unmarshal`对（对应`Create`和`Parse`）  
也就是针对不同的消息种类我们的`Create`和`Parse`函数都不同，而若我们将所有类型的消息都定义在一个结构体内，发任何消息都直接Marshal，Unmarshal即可

# 21.Jul
1. 使用json文件初始化参数运行，并通过flag包用命令行指定配置文件：
```
go run main.go -config config.json
```
2. 完成了20.Jul 要求的统一序列化和反序列化，网络传输只使用了一个结构 `NetMsg` 进行通信，这样序列化`Create`和反序列化`Parse`就只需要使用`json.Marshal`和`json.Unmarshal`就可以了
```go
// protocol/typ.go
type NetMsg struct {
	Id           string
	Pswd         string
	Msg          string
	RoomId       string
	LoginSuccess bool
}
```
统一序列化和反序列化放在了 `UniMessage`分支下面