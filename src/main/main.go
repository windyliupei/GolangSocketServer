package main

//import protocol "protocol"

import (
	"fmt"
	"io"
	"net"
	"os"
	"protocol"
	"time"
	"xmlConfig"

	"bufio"
	"log"
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

	running := true
	reader := bufio.NewReader(os.Stdin)
	for running {
		data, _, _ := reader.ReadLine()
		command := string(data)
		if command == "stop" {
			running = false
		}
		log.Println("command", command)
	}
	if !running {
		os.Exit(1)
	}

}

//处理连接
func handleConnection(conn net.Conn) {

	// 缓冲区，存储被截断的数据
	tmpBuffer := make([]byte, 0)

	//接收解包
	readerChannel := make(chan []byte, 16)
	go reader(readerChannel)

	buffer := make([]byte, 1024) //TODO:Need handle the buffer size.
	//for {
	n, err := conn.Read(buffer)

	if err != nil {
		if err != io.EOF {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
		}
		return
	}

	tmpBuffer = protocol.Depack(append(tmpBuffer, buffer[:n]...), readerChannel)
	//Return to the client side.
	conn.Write(tmpBuffer)

	//timeout := 10
	//messnager := make(chan byte)
	//心跳计时
	//go HeartBeating(conn, messnager, timeout)
	//检测每次Client是否有数据传来
	//Log("!" + string(tmpBuffer))
	//go GravelChannel(tmpBuffer, messnager)
	//}

	defer conn.Close()

}

func reader(readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			Log("receive data string:" + string(data))
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

//心跳计时，根据GravelChannel判断Client是否在设定时间内发来信息
func HeartBeating(conn net.Conn, messnager chan byte, timeout int) {
	select {
	case fk := <-messnager:
		Log(conn.RemoteAddr().String(), "receive data string:", string(fk))
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		break
	case <-time.After(time.Second * 5):
		Log("It's really weird to get Nothing!!!")
		conn.Close()
	}

}

func GravelChannel(n []byte, messnager chan byte) {

	Log("~!" + string(n))
	Log(len(n))
	for _, v := range n {
		Log(string(v))
		messnager <- v
	}
	close(messnager)
}
