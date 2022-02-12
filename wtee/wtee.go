package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/ahuigo/glogger"
)

var logger = glogger.Glogger

//client
func main() {
	var err error
	// nc -l 4600
	hostPortPtr := flag.String("h", "127.0.0.1:4600", "write host")
	flag.Parse()

	tails := flag.Args()
	hostPortSlice := strings.Split(*hostPortPtr, ":")
	host := hostPortSlice[0]
	port := "4600"
	if len(hostPortSlice) >= 2 {
		port = hostPortSlice[1]
	}
	conn, _ := createConn(host + ":" + port)
	fileReader := os.Stdin
	if len(tails) > 0 {
		if tails[0] != "-" {
			if fileReader, err = os.Open(tails[0]); err != nil {
				perror("Invalid path: ", tails[0], err)
				// os.Exit(1)
			}
		}
	}

	logger.Info(host, port, fileReader)
	// reader to tcp writer
	for {
		buf := make([]byte, 1000)
		n, err := fileReader.Read(buf)
		if err != nil {
			perror(err)
		}
		if n > 0 {
			writeToServer(conn, buf[:n])
			// fileWriter.Write(buf)
		} else {
			break
		}
	}
}

func createConn(addr string) (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", addr)
	if err != nil {
		perror("Connect to TCP server failed ,err:", err)
		return
	}
	return
}

func writeToServer(conn net.Conn, buf []byte) {
	_, err := conn.Write(buf)
	if err != nil {
		perror("Write failed,err:", err)
	}
}

func perror(args ...interface{}) {
	// fmt.Println(args...)
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}
