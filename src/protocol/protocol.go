//通讯协议处理
package protocol

import (
	"bytes"
	"encoding/binary"
	"strings"
)

const (
	ConstHeader       = "Headers"
	ConstHeaderLength = 7
	ConstMLength      = 4
)

//封包
func Enpack(message []byte) []byte {
	return append(append([]byte(ConstHeader), IntToBytes(len(message))...), message...)
}

//解包
/*
func Depack(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+ConstHeaderLength+ConstMLength {
			break
		}
		if string(buffer[i:i+ConstHeaderLength]) == ConstHeader {
			messageLength := BytesToInt(buffer[i+ConstHeaderLength : i+ConstHeaderLength+ConstMLength])
			fmt.Println(messageLength)
			if length < i+ConstHeaderLength+ConstMLength+messageLength {
				break
			}
			data := buffer[i+ConstHeaderLength+ConstMLength : i+ConstHeaderLength+ConstMLength+messageLength]
			readerChannel <- data

		}
	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}
*/

func Depack(buffer []byte, readerChannel chan []byte) []byte {

	length := len(buffer)
	var i int

	if length < 2 { //each transfer body format should be [KEY]:[messageText]#,so it will not less than 2 length
		return make([]byte, 0)
	}

	var transferBody = string(buffer[i:])
	var prefix = substr(transferBody, 0, strings.Index(transferBody, ":"))

	switch key := prefix; key {
	case "echo":
		{
			var messageBody = strings.TrimSuffix(strings.TrimPrefix(transferBody, key+":"), "#")
			data := []byte(messageBody)
			readerChannel <- data
			return data
		}

	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

//截取字符串 start 起点下标 end 终点下标(不包括)
func substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}
