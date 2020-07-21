package appsocket

import (
	"chat_v3/protocol"
	"log"
	"net"
)

const BUFFER_SIZE = 1024 // conn.Read函数所使用的buffer的大小
const TYPE_SIZE = 4
const LENGTH_SIZE = 4

type AppSocket struct {
	conn   net.Conn
	buf    []byte
	nBytes uint32 // buf中总共n位是有效的
	flag   uint32
}

func NewAppSocket(conn net.Conn) *AppSocket { // 这里要在函数内声明变量然后返回地址
	var as AppSocket
	as.conn = conn
	as.buf = make([]byte, BUFFER_SIZE)
	as.nBytes = 0
	as.flag = 0
	return &as
}

func (as *AppSocket) Close() {
	as.conn.Close()
}

// ReadAppFrame 从conn中读取一个TLV结构
func (as *AppSocket) ReadAppFrame() (ft protocol.FrameType, val []byte, err error) {

	// 从conn中读取4个字节作为TYPE
	typeBytes, err := as.readLenN(TYPE_SIZE)
	//fmt.Println("typeBytes len =", len(typeBytes))
	if err != nil {
		return ft, val, err
	}
	ft = protocol.FrameType(protocol.BytesToInt(typeBytes))

	// 从conn中读取4个字节作为LENGTH
	lenBytes, err := as.readLenN(LENGTH_SIZE)
	//fmt.Println("lenBytes len =", len(lenBytes))
	if err != nil {
		return ft, val, err
	}
	length := protocol.BytesToInt(lenBytes)
	//log.Println("MSG LENGTH =", length)

	// 读取LENGTH个字节作为VAL
	val, err = as.readLenN(length)
	//fmt.Println("valBytes = ", len(val))
	if err != nil {
		return ft, val, err
	}

	return ft, val, err
}

func (as *AppSocket) WriteAppFrame(content []byte) (n int, err error) {
	n, err = as.conn.Write(content)
	return n, err
}

// mvBytes2Front 用来把n个字节的数据提到buf的最前面, 成功返回true
func mvBytes2Front(buf []byte, from, n int) bool {
	if from+n > len(buf) {
		log.Println("复制长度超出buffer长度")
		return false
	}
	copy(buf[0:n], buf[from:from+n])
	return true
}

// readLenN 从conn中读取len个byte
func (as *AppSocket) readLenN(len uint32) (val []byte, err error) {
	if len == 0 {
		return val, nil
	}
	left := uint32(as.nBytes - as.flag) // buf 中还有left个byte未读取
	for len > 0 {                       // 当还未读满时
		//fmt.Println("readLenN 嘟嘟嘟嘟")
		if len <= left { // 剩下未读取byte数量大于所需数量，无需再从conn中收取新的
			val = append(val, as.buf[as.flag:as.flag+len]...)
			as.flag += len // 读取了len个byte后更新flag
			return val, nil
		} else { // buf中的byte数量不够，需要从conn中读取新的数据
			val = append(val, as.buf[as.flag:as.nBytes]...) // 先把buf里剩下的全读完
			readBytes := as.nBytes - as.flag                // 这次读取了 readBytes 个 byte
			len = len - readBytes                           // 还需要读取 needBytes 个 byte
			n, err := as.conn.Read(as.buf)                  // 从conn中收取
			if err != nil {
				return val, err
			}
			// 读取后重设 flag 和 nBytes 参数
			as.nBytes = uint32(n)
			as.flag = 0
			left = as.nBytes - as.flag
		}
	}
	return val, err
}
