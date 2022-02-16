//go:build rtee
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ahuigo/wrtee/file"
)

type RfStatus int

const bodyDilimiter = '\n'
const (
	RfStatusInit RfStatus = iota + 1
	RfStatusRecv
)

type RecvFiles struct {
	// RfStatusInit: check header ; parse header+ create file
	// RfStatusRecv: recv data ; EOF: close file + init
	status RfStatus
	fp     *os.File
	home   string
	force  bool
}

type Args struct {
	Port  string
	Home  string //output home
	Force bool
}

func getCliArgs() (args Args) {
	portPtr := flag.String("p", "8100", "read port")
	forcePtr := flag.Bool("f", false, "force sync")
	flag.Parse()
	// port
	args.Port = *portPtr
	args.Force = *forcePtr
	debug("force:", *forcePtr)

	// dir
	tailArgs := flag.Args()
	if len(tailArgs) > 0 {
		args.Home = tailArgs[0]
	} else {
		args.Home, _ = os.Getwd()
	}
	return args
}

func (rf *RecvFiles) ReadHeader(b *[]byte) (finished bool, err error) {
	buf := *b
	sepIndex := bytes.IndexByte(buf, bodyDilimiter)
	if sepIndex > 0 {
		fmt.Println("ahui: parse header")
		header := string(buf[:sepIndex])
		segs := strings.Split(header, ":")
		if len(segs) < 2 || segs[0] != "file" {
			return false, errors.New("invalid header:" + header)
		}
		// finished
		fpath := filepath.Join(rf.home, segs[1])
		rf.fp, err = file.CreateFile(fpath, rf.force)
		*b = (*b)[sepIndex+1:]

		return true, err
	}
	return false, nil
}

func (rf *RecvFiles) ReadFile(b *[]byte) (finished bool, readNext bool, err error) {
	buf := *b
	sepIndex := bytes.IndexByte(buf, ':')
	fmt.Println("read file with string:", string(buf))
	if sepIndex > 0 {
		segName := string(buf[:sepIndex])
		segLen, err := strconv.Atoi(segName)
		if segName == "END" {
			// end
			*b = (*b)[sepIndex+1:]
			rf.fp.Close()
			finished = true
			return finished, false, nil
		} else if err == nil && segLen > 0 {
			finished = false
			// write line
			if len(buf[sepIndex+1:]) >= segLen {
				// non-block
				bytes := buf[sepIndex+1 : sepIndex+1+segLen]
				file.WriteBytes(rf.fp, bytes)
				*b = (*b)[sepIndex+1+segLen:]
				return finished, false, nil
			} else {
				// block line
				readNext = true
				return finished, true, nil
			}
		} else {
			// bad seg length
			return finished, readNext, errors.New("bad body segName:" + segName)
		}

	} else {
		readNext = true
		if len(buf) > 10 { // len(number_len|END:)<10
			return false, false, errors.New("bad body seg:" + string(buf))

		}
	}
	return
}

/**
 * file:file_path1
 * \n
 * length:data
 * EOF:
 * file:file_path2
 * \n
 * length:data
 * EOF:
 */
func (rf *RecvFiles) Read(b *[]byte) (err error) {
	finished := false
	readNext := false //read more byte(not enough bytes)
	for {
		switch rf.status {
		case RfStatusInit:
			fmt.Println("parse header...:", string(*b))
			if finished, err = rf.ReadHeader(b); err != nil {
				return err
			}
			if finished {
				rf.status = RfStatusRecv
			} else {
				readNext = true
			}
		case RfStatusRecv:
			fmt.Println("parse file...:", string(*b))
			if finished, readNext, err = rf.ReadFile(b); err != nil {
				return err
			}
			if finished {
				rf.status = RfStatusInit
			} else {
				readNext = true
			}
		default:
			return errors.New("undefined recv file status")
		}
		if readNext {
			break
		}

	}
	return
}

func recvConn(conn net.Conn, rf *RecvFiles) {
	defer conn.Close()
	var buf [5]byte
	bytes := make([]byte, 0, 1000)
	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			if err != io.EOF {
				fmt.Println("Read from tcp server failed,err:", err)
			} else {
				fmt.Print("Read from tcp server EOF")
			}
			break
		}
		data := string(buf[:n])
		fmt.Printf("Recived from client:%s\n", data)
		bytes = append(bytes, buf[:n]...)
		fmt.Printf("Recived bytes:%s\n", string(bytes))
		if err = rf.Read(&bytes); err != nil {
			fmt.Printf("read err:%+v\n", err)
			conn.Close()
			break
		}

	}
}

func main() {
	// 监听TCP 服务端口
	args := getCliArgs()
	listener, err := net.Listen("tcp", "0.0.0.0:"+args.Port)
	if err != nil {
		fmt.Println("Listen tcp server failed,err:", err)
		return
	}
	fmt.Println("Listen port:", args.Port)
	fmt.Println("Home:", args.Home)

	for {
		// 建立socket连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listen.Accept failed,err:", err)
			continue
		}

		// 业务处理逻辑
		rf := &RecvFiles{
			status: RfStatusInit,
			home:   args.Home,
			force:  args.Force,
		}
		debug("rf:", rf)
		go recvConn(conn, rf)
	}
}

func debug(ds ...interface{}) {
	for _, d := range ds {
		fmt.Printf("%#v ", d)
	}
	fmt.Printf("\n")
}
