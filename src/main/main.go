package main 

import protocol "protocol"

import (
"fmt"  
    "net"  
    "os"  
   )

func main() {
	//建立socket，监听端口  
    netListen, err := net.Listen("tcp", "localhost:1024")  
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
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())    
        os.Exit(1)    
    }    
}    