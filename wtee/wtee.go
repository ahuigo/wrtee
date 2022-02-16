package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/ahuigo/glogger"
	"github.com/ahuigo/wrtee/file"
)

var logger = glogger.Glogger

type Client struct {
	addr  string
	paths []string //output home
	conn  net.Conn
}

func (c *Client) getConn() net.Conn {
	var err error
	if c.conn == nil {
		c.conn, err = net.Dial("tcp", c.addr)
		if err != nil {
			perror("Connect to TCP server failed ,err:", err)
		}
	}
	return c.conn
}
func (c *Client) sendPaths() (err error) {
	for _, path := range c.paths {
		err = c.sendPath(path)
		if err != nil {
			return
		}
	}
	return
}
func (c *Client) sendPath(path string) (err error) {
	if file.IsDir(path) {

	} else {
		err = c.sendFile(path)
	}
	return
}
func (c *Client) sendFile(filePath string) (err error) {
	fileReader, _ := c.openFile(filePath)
	segLenBits := 5
	buf := make([]byte, 2000)

	header := fmt.Sprintf("file:%s\n", file.GetFilename(filePath))
	c.sendBlob([]byte(header))
	defer func() {
		c.sendBlob([]byte("END:"))
	}()
	for {
		n, err := fileReader.Read(buf[segLenBits+1:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if n > 0 {
			segLen := fmt.Sprintf("%05d:", n)
			copy(buf, []byte(segLen))
			c.sendBlob(buf[:n+segLenBits+1])
			// fileWriter.Write(buf)
		} else {
			break
		}
	}
	return
}

func (c *Client) sendBlob(buf []byte) {
	fmt.Println("sendBlob:", string(buf))
	conn := c.getConn()
	n, err := conn.Write(buf)
	if err != nil {
		perror("Write failed,err:", err)
	}else{
        if n!=len(buf){
            perror("not valid n:", n, "len=",len(buf))
        }
    }
}

func (c *Client) openFile(filePath string) (fileReader io.Reader, err error) {
	if filePath == "-" {
		fileReader = os.Stdin
	} else {
		if fileReader, err = os.Open(filePath); err != nil {
			perror("Invalid path: ", filePath, err)
			// os.Exit(1)
		}
	}
	return
}

//     go run wtee.go -h 127.0.0.1:4600 a.txt
func getCliArgs() (args Client) {
	addrPtr := flag.String("h", "127.0.0.1:8100", "send server")
	flag.Parse()
	// port
	args.addr = *addrPtr

	// dir
	args.paths = flag.Args()
	if len(args.paths) == 0 {
		args.paths = []string{"-"}
	}
	// conn, _ := createConn(host + ":" + port)
	return args
}

//client
func main() {
	client := getCliArgs()
	client.sendPaths()
}

func perror(args ...interface{}) {
	// fmt.Println(args...)
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}
