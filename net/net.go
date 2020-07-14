package net

import (
	"log"
	"net"
)

const BUFFER_SIZE = 1024  // conn.Read函数所使用的buffer的大小
const TAG_SIZE = 4
const LENGTH_SIZE = 4


// 这里要写一个application layer的读和写函数
// 从网络中读一个应用层TLV(Type Length Value)数据包并返回V

type NetReadWriter struct {
	conn net.Conn
	buf []byte
	flag uint32   // 读取到的位置
	nBytes uint32 // buf中总共n位是有效的
}
func NewNetReaderWriter(conn net.Conn) *NetReadWriter{  // 这里要在函数内声明变量然后返回地址
	var nrw NetReadWriter
	nrw.conn = conn
	nrw.buf = make([]byte, BUFFER_SIZE)
	nrw.nBytes = 0
	nrw.flag = 0
	return &nrw
}

// ReadAppFrame 从conn中读取一个TLV结构
func (nrw *NetReadWriter) ReadAppFrame(content []byte) (tag ProtocalTag, val []byte){
	var len uint32                // TLV 中V的长度
	tag, len = nrw.findTag()      // 找到TAG 和 LEN
	nrw.readLenN(val, len)        // 读取len个byte
	return tag, val
}

func (nrw *NetReadWriter) WriteAppFrame(content []byte) (n int, err error){
	n, err = nrw.conn.Write(content)
	return n, err
}

// mvBytes2Front 用来把n个字节的数据提到buf的最前面, 成功返回true
func mvBytes2Front(buf []byte, from, n int) bool {
	if from + n > len(buf) {
		log.Println("复制长度超出buffer长度")
		return false
	}
	copy(buf[0:n], buf[from: from+n])
	return true
}

// findTag 从conn中收取数据，找到tag ， 返回 tag 和 数据包长度
func (nrw *NetReadWriter)findTag() (tag ProtocalTag, len uint32){
	for ;nrw.flag <= nrw.nBytes - 4; nrw.flag++ {

	}

}

// readLenN 从conn中读取len个byte
func (nrw *NetReadWriter)readLenN(val []byte, len uint32){
	if len == 0{
		return
	}
	left := uint32(nrw.nBytes - nrw.flag)   // buf 中还有left个byte未读取
	if len < left {                         // 无需再从conn中收取新的
		val = append(val, nrw.buf[nrw.flag:nrw.flag+len]...)
		nrw.flag += len                     // 读取了len个byte后更新flag
		return
	} else {                                // buf中的byte数量不够，需要从conn中读取
		val = append(val, nrw.buf[nrw.flag:nrw.nBytes]...)
		readBytes := nrw.nBytes - nrw.flag  // 这次读取了 readBytes 个 byte
		needBytes := len - readBytes        // 还需要读取 needBytes 个 byte
		n, _ := nrw.conn.Read(nrw.buf)      // 从conn中收取
		// 读取后重设 nrw 中的 flag 和 nBytes 参数
		nrw.nBytes = uint32(n)
		nrw.flag = 0
		nrw.readLenN(val, needBytes)        // 再次读取，直到完全读完
	}
}

