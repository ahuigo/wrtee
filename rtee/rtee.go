package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func process(conn net.Conn) {
	defer conn.Close()
	for {
		var buf [5]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			perror("Read from tcp server failed,err:", err)
			break
		}
		data := string(buf[:n])
		os.Stdout.Write(buf[:])

		fmt.Fprintf(os.Stderr, "Recived from client,data:%s,len=%d\n", data, len(data))
	}
}

func getPort() (port string) {
	portPtr := flag.String("p", "4600", "read port")
	flag.Parse()
	port = *portPtr
	return
}

func main() {
	// 监听TCP 服务端口
	port := getPort()
	listener, err := net.Listen("tcp", "0:"+port)
	if err != nil {
		perror("Listen tcp server failed,err:", err)
		return
	}
	println("listen port:", port)

	for {
		// 建立socket连接
		conn, err := listener.Accept()
		if err != nil {
			perror("Listen.Accept failed,err:", err)
			continue
		}

		// 业务处理逻辑
		go process(conn)
	}
}

func perror(args ...interface{}) {
	// fmt.Println(args...)
	fmt.Fprintln(os.Stderr, args...)
	// os.Exit(1)
}
