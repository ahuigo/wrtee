//go:build !rtee
package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/ahuigo/wrtee/file"
)

type Args struct {
	addr string //127.0.0.1:4600
	conn net.Conn
	srcs []string //send path
	dst  string   //recv path
}

//     zsync -s 127.0.0.1:4600/path source[source2 ... sourceN] destination
func getClient() (args Args) {
	srvHostPortPtr := flag.String("s", "127.0.0.1:4600/", "server url")
	flag.Parse()

	args.addr = strings.Split(*srvHostPortPtr, "/")[0]
	args.srcs = flag.Args()
	if len(args.srcs) == 0 {
		args.srcs = []string{"-"}
	}

	fmt.Println("args:", args)
	// args.connect()

	return args
}

func (c *Args) connect() (err error) {
	c.conn, err = net.Dial("tcp", c.addr)
	if err != nil {
		panic(err)
	}
	return err
}

func (c *Args) syncSrcs() (err error) {
	for _, src := range c.srcs {
		err = c.syncSrc(src)
	}
	return
}
func (c *Args) syncSrc(src string) (err error) {
	if src == "-" {
		c.syncFile("stdin.txt", os.Stdin)
	} else if file.IsDir(src) {

	} else {
		// singfile
		if srcReader, err := os.Open(src); err != nil {
			perror("Invalid path: ", src, err)
			// os.Exit(1)
		} else {
			err = c.syncFile(file.GetFilename(src), srcReader)
		}
	}
	return
}

// func (c *Args) syncFile(src string, fp *os.File) (err error) {
func (c *Args) syncFile(src string, fp io.Reader) (err error) {
	// _, err = io.Copy(os.Stdout, fp) // copy b to stdout
	// _, err = io.Copy(c.conn, fp) // copy b to stdout
	c.write([]byte("file:" + src + "\n"))
	segLenBits := 6
	bbuf := make([]byte, 100)
	for {
		// buf := new(bytes.Buffer)
		// // buf.ReadFrom(fp)
		// buf.Read(bbuf)

		n, err := fp.Read(bbuf[segLenBits:])
		if err != nil {
			if err == io.EOF {
				err = c.write([]byte("END:"))
				return nil
			}
			return err
		}

		// write length
		copy(bbuf, []byte(fmt.Sprintf("%05d:", n)))
		// nBits, err := lenReader.Read(bbuf[:segLenBits])

		err = c.write(bbuf[:segLenBits+n])
		if err != nil {
			return err
		}
	}
}
func (c *Args) write(b []byte) (err error) {
	fmt.Println("send buf seg:", string(b), "len=", len(b), "b=", b)
	if true {
		return nil
	}
	if c.conn == nil {
		c.connect()
	}
	_, err = c.conn.Write(b)
	return
}

//client
func main() {
	client := getClient()
	client.syncSrcs()

}

func perror(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(1)
}
