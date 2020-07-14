package net

import (
	"bytes"
	"encoding/binary"
)


const TAG_HEAD_0 byte = 0xFF
const TAG_HEAD_1 byte = 0xFE
const TAG_HEAD_2 byte = 0xFD


type ProtocalTag uint8
const(
	TAG_LOGIN ProtocalTag = 101  // 登录
	TAG_LOGOUT ProtocalTag = 102 // 注销
	TAG_JOIN ProtocalTag = 103   // 加入聊天室
	TAG_QUERY ProtocalTag = 104  // 查询聊天室
	TAG_MSG ProtocalTag = 105    // 消息
	TAG_LOGIN_FAIL ProtocalTag = 201 // 登录失败
	TAG_ROOMS ProtocalTag = 202  // 聊天室列表
	NO_TAG ProtocalTag = 001
	INCOMPLETE_TAG ProtocalTag = 002
)

// IsTag 判断是否是Tag
func IsTag(tag [4]byte) (ProtocalTag, bool) {
	//if tag[0] == TAG_HEAD && tag[1] == TAG_HEAD && tag[2] == TAG_HEAD {
	//		switch tag[3]{
	//		case 101:
	//			return TAG_LOGIN, true
	//		case 102:
	//			return TAG_LOGOUT, true
	//		case 103:
	//			return TAG_JOIN, true
	//		case 104:
	//			return TAG_QUERY, true
	//		case 105:
	//			return TAG_MSG, true
	//		case 201:
	//			return TAG_LOGIN_FAIL, true
	//		case 202:
	//			return TAG_ROOMS, true
	//		}
	//}
	if tag[0] == TAG_HEAD_0 && tag[1] == TAG_HEAD_1 && tag[2] == TAG_HEAD_2 && (
		tag[3] == 101 || tag[3] == 102 || tag[3] == 103 || tag[3] == 104 || tag[3] == 105 || tag[3] == 201 || tag[3] == 202){
			return ProtocalTag(tag[3]), true
	}
	return NO_TAG , false
}

func IsIncompleteTag(tag []byte) {
	len := len(tag)

}


// 下面的函数根据内容和创建TLV数据包
func CreateLoginBS(id, pswd string) []byte {

}

func CreateJoinBS(roomid string) []byte {

}

func CreateQueryBS() []byte{

}

func getTagBS(tag ProtocalTag) []byte {
	bs := make([]byte, 4)
	bs[0] = TAG_HEAD_0
	bs[1] = TAG_HEAD_1
	bs[2] = TAG_HEAD_2
	bs[3] = byte(tag)
	return bs
}

func getLengthBS(len uint32)[]byte{
	return intToBytes(len)
}

// intToBytes 将1个32位无符号整数转化为4个byte
func intToBytes(n uint32) []byte {
	x := uint32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ := binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// bytesToInt 将4个byte转化为一个32位无符号整数
func bytesToInt(b []byte) uint32 {
	bytesBuffer := bytes.NewBuffer(b)
	var x uint32
	_ := binary.Read(bytesBuffer, binary.BigEndian, &x)
	return uint32(x)
}

