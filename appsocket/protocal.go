package appsocket

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"strings"
)


const TAG_HEAD_0 byte = 0xFF
const TAG_HEAD_1 byte = 0xFE
const TAG_HEAD_2 byte = 0xFD


type FrameType uint32
const(
	T_LOGIN FrameType = 101  // 登录
	T_LOGIN_SUCCESS FrameType = 102 // 登录成功
	T_JOIN FrameType = 103   // 加入聊天室
	T_QUERY FrameType = 104  // 查询聊天室
	T_MSG FrameType = 105    // 消息
	T_LOGIN_FAIL FrameType = 201 // 登录失败
	T_ROOMS FrameType = 202  // 聊天室列表
)

type LogInfo struct {
	id string
	pswd string
}

// 下面的函数根据内容和创建TLV数据包
func CreateLoginBS(id, pswd string) []byte {
	loginfo := id + " " + pswd
	return CreateBS(T_LOGIN, loginfo)

}

func CreateJoinBS(roomid string) []byte {
	return CreateBS(T_JOIN, roomid)
}

func CreateQueryBS() []byte{
	return CreateBS(T_QUERY, "")
}

func CreateMsgBS(msg string) []byte{
	return CreateBS(T_MSG, msg)
}

func CreateBS(ft FrameType, value string) []byte{
	var bs []byte
	T := typeToBytes(ft)
	V := []byte(value)
	L := lenToBytes(uint32(len(V)))
	bs = append(bs, T...)
	bs = append(bs, L...)
	bs = append(bs, V...)
	return bs
}

// ParseLogInfo 从字节流解析获得用户名、密码信息
func ParseLogInfo(bs []byte) (id, pswd string){
	str := string(bs)
	id = strings.Fields(str)[0]
	pswd = strings.Fields(str)[1]
	fmt.Println("解析得到用户名密码为：", id, pswd)
	return id, pswd
}

// ParseRoomId 从字节流解析获得房间名称
func ParseRoomId(bs []byte) (r string) {
	return string(bs)
}

// 以下转换均为大端模式（低地址放高位）
// i := 0x11         22         33         44
// b[0] = 11, b[1] = 22, b[2] = 33, b[3] = 44

// bytesToType 将4个byte转为一个FrameType
func bytesToType(bs []byte) FrameType{
	return FrameType(bytesToInt(bs))
}

// bytesToType 将4个byte转为一个uint32类型变量
func bytesToLen(bs []byte) uint32{
	return uint32(bytesToInt(bs))
}

// typeToBytes 将一个FrameType转为4个byte
func typeToBytes(t FrameType) []byte {
	return intToBytes(uint32(t))
}

// lenToBytes 将一个uint32类型转为4个byte
func lenToBytes(len uint32) []byte {
	return intToBytes(len)
}

// intToBytes 将1个32位无符号整数转化为4个byte
func intToBytes(n uint32) []byte {
	x := uint32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	err := binary.Write(bytesBuffer, binary.BigEndian, x)
	if err != nil {
		log.Fatal("intToBytes", err)
	}
	return bytesBuffer.Bytes()
}

// bytesToInt 将4个byte转化为一个32位无符号整数
func bytesToInt(b []byte) uint32 {
	bytesBuffer := bytes.NewBuffer(b)
	var x uint32
	err := binary.Read(bytesBuffer, binary.BigEndian, &x)
	if err != nil {
		log.Fatal("bytesToInt", err)
	}
	return uint32(x)
}

