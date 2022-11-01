package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

// handler 应用服务器的抽象
type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}

// 监听并提供服务，并在收到 closeChan 发来的关闭通知后关闭
func ListenAndServe(address string, closeCh <-chan struct{}) {
	log.Println("listen starting ...")
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		log.Println("listener closing2 ...:")
		listener.Close()
		log.Println("listener closed2!")
	}()



	// 监听通知
	go func(closeCh <-chan struct{}) {
		<-closeCh
		log.Println("listener closing ...:")
		listener.Close()
		log.Println("listener closed!")
	}(closeCh)

	for {
		fmt.Println("listener accepting ...")
		// 监听请求
		conn, err := listener.Accept()
		if err != nil {
			log.Panicln(err)
			continue
		}
		fmt.Println("listener accepted!")
		// 处理请求
		go Handle(conn)

	}
}

func Handle(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("发送了panic")
		}
	}()

	fmt.Println("listener handing ...")
	r := bufio.NewReader(conn)

	//panic(8)

	for {
		data, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("listener exit!")
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
	ch := make(chan struct{})
	go ListenAndServe(":8080", ch)
	time.Sleep(15 * time.Second)
	fmt.Println("发送关闭listener信号 ...")
	ch <- struct{}{}
}
