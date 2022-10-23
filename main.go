package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

func ListenAndServe(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	log.Println("listen:", address)
	for {
		// 监听请求
		conn, err := listener.Accept()
		if err != nil {
			log.Panicln(err)
			continue
		}
		// 处理请求
		go Handler(conn)
	}
}

func Handler(conn net.Conn) {
	r := bufio.NewReader(conn)

	for {
		data, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("conn :%s closed.", conn.LocalAddr().String())
			} else {
				log.Println("Handler err:", err)
			}
			return
		}
		// 将收到的信息发送给客户端
		conn.Write(data)
	}

}

func main() {
	ListenAndServe(":8080")
}
