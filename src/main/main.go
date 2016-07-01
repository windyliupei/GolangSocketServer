package main

//import protocol "protocol"

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"../protocol"
)

func main() {

	var configMap map[string]string
	configMap = getConfig()

	localhost := getElement("host", configMap)
	port := getElement("port", configMap)

	//建立socket，监听端口
	netListen, err := net.Listen("tcp", localhost+":"+port)
	CheckError(err)
	defer netListen.Close()

	Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		Log(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(conn)
	}
}

//处理连接
func handleConnection(conn net.Conn) {

	// 缓冲区，存储被截断的数据
	tmpBuffer := make([]byte, 0)

	//接收解包
	readerChannel := make(chan []byte, 16)
	go reader(readerChannel)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		tmpBuffer = protocol.Depack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}
	defer conn.Close()

}

func reader(readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			Log(string(data))
		}
	}
}

func Log(v ...interface{}) {
	fmt.Println(v...)
}

func CheckError(err error) {
	if err != nil {
		var errText = err.Error()
		fmt.Println(errText)
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

//Get Host address and prot frm configuration file.
func getConfig() map[string]string {

	var token xml.Token
	var err error

	configMap := make(map[string]string)

	currentPath, err := os.Getwd()
	CheckError(err)

	configByte, err := ioutil.ReadFile(currentPath + "/config.xml")

	decoder := xml.NewDecoder(bytes.NewBuffer(configByte))

	var elementCount int

	var keyList []string
	var valueList []string

	for token, err = decoder.Token(); err == nil; token, err = decoder.Token() {

		switch token := token.(type) {
		// 处理元素
		case xml.StartElement:
			name := token.Name.Local
			keyList = append(keyList, name)
			fmt.Println(keyList)
			elementCount++
		case xml.CharData:
			content := string([]byte(token))
			if content != "\n" {
				valueList = append(valueList, content)
				fmt.Println(valueList)
			}

		}
	}
	for index := 0; index < elementCount; index++ {
		configMap[keyList[index]] = valueList[index]
	}
	fmt.Println(elementCount)

	return configMap
}

func getElement(key string, themap map[string]string) string {
	if value, ok := themap[key]; ok {
		return value
	}

	return ""
}
