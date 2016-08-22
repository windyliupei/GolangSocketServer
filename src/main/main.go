package main

//import protocol "protocol"

import (
	"fmt"
	"io"
	"net"
	"os"
	"xmlConfig"

	"protocol"
)

func main() {

	fmt.Println("Server Start")
	//Database access : get mac encrypt key.
	//dbAccess.GetMacEncryptKey("00d02d23d2ae")

	//Get Socket setting, host, port from xml file.
	var configMap map[string]string
	configMap = xmlConfig.GetConfig()

	localhost := xmlConfig.GetElement("host", configMap)
	port := xmlConfig.GetElement("port", configMap)

	//建立socket，监听端口
	netListen, err := net.Listen("tcp", localhost+":"+port)
	Log("Localhost:" + localhost + "port:" + port)
	CheckError(err)
	defer netListen.Close()

	Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		Log(conn.RemoteAddr().String(), " tcp connect successful!")
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
		Log(n)
		if err != nil {
			if err != io.EOF {
				Log(conn.RemoteAddr().String(), " connection error: ", err)
			}
			return
		}

		tmpBuffer = protocol.Depack(append(tmpBuffer, buffer[:n]...), readerChannel)

		conn.Write(tmpBuffer)
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
