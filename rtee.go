//go:build !rtee

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
	"github.com/ahuigo/wrtee/util"
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
	status         RfStatus
	fp             *os.File
	home           string
	force          bool
	recvOffset     int64
	recvRealOffset int64
	conn           net.Conn
}

type Args struct {
	Port  string
	Home  string // save file at home
	Force bool
}

func getCliArgs() (args Args) {
	portPtr := flag.String("p", "8100", "read port")
	forcePtr := flag.Bool("f", false, "force sync")
	flag.Parse()
	// port
	args.Port = *portPtr
	args.Force = *forcePtr
	println("force:", *forcePtr)

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
		// util.Debugf("ahui: parse header")
		header := string(buf[:sepIndex])
		segs := strings.Split(header, ":")
		if len(segs) < 2 || segs[0] != "file" {
			return false, errors.New("invalid header:" + header)
		}
		// finished
		fpath := filepath.Join(rf.home, segs[1])
		rf.fp, err = file.CreateFile(fpath, rf.force)
		*b = (*b)[sepIndex+1:]
		fmt.Printf("receive file:%s\n", fpath)

		return true, err
	}
	return false, nil
}

var output []byte

func (rf *RecvFiles) Read2File(b *[]byte) (finished bool, isReadMore bool, err error) {
	buf := *b
	sepIndex := bytes.IndexByte(buf, ':')
	if sepIndex > 0 {
		segName := string(buf[:sepIndex])
		segLen, err := strconv.Atoi(segName)
		if segName == "END" {
			// end
			*b = (*b)[sepIndex+1:]
			rf.fp.Close()
			finished = true
			return finished, isReadMore, nil
		} else if err == nil && segLen > 0 {
			finished = false
			// write line
			if len(buf[sepIndex+1:]) >= segLen {
				// non-block
				bytes := buf[sepIndex+1 : sepIndex+1+segLen]
				err := file.WriteBytes(rf.fp, bytes)
				*b = (*b)[sepIndex+1+segLen:]

				// ********************debug ***************************
				if false {
					rf.debugRead2File(bytes)
				}
				// ********************debug end***************************

				return finished, isReadMore, err
			} else {
				// block line
				isReadMore = true
				return finished, isReadMore, nil
			}
		} else {
			// bad seg length
			return finished, isReadMore, errors.New("bad body segName:" + segName)
		}

	} else {
		isReadMore = true
		if len(buf) > 10 { // len(number_len|END:)<10
			return false, false, errors.New("bad body seg:" + string(buf))

		}
	}
	return
}

func (rf *RecvFiles) debugRead2File(bytes []byte) {
	byten := int64(len(bytes))
	oribytes, err := file.ReadSeek("go2", rf.recvOffset, byten)
	if err != nil {
		util.Fatalf("seek read err:%s,offset:%d,n:%d", err, rf.recvOffset, byten)
	}
	offset := rf.recvOffset
	if dn := file.BytesDiffn(oribytes, bytes); dn >= 0 {
		util.Debugf("recv diff at offset:%d,n:%d\n bytes:%v\n oribytes:\n%v", offset, dn, string(bytes), string(oribytes))
		// util.Debugf("dif2\n bytes:%v\nobytes=%v", bytes, oribytes)
		os.Exit(1)
	} else {
		// util.Debugf("save offset %d(%d):%v", offset+byten, byten, string(bytes))
		util.Debugf("offset %d ok", offset+byten)
		// util.Debugf("offset %d ok:\n%s", offset, string(bytes))
		output = append(output, bytes...)
		// if len(output) > 3888 {
		// 	util.Debugf("output(%d):\n%s", offset, string(output))
		// 	os.Exit(1)
		// }
	}
	rf.recvOffset += byten
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
func (rf *RecvFiles) handleRecvBytes(b *[]byte) (err error) {
	finished := false
	isReadMore := false //read more byte(not enough bytes)
	for {
		switch rf.status {
		case RfStatusInit:
			if finished, err = rf.ReadHeader(b); err != nil {
				return err
			}
			if finished {
				rf.status = RfStatusRecv
			} else {
				isReadMore = true
			}
		case RfStatusRecv:
			// util.Debugf("parse file...:%s", string(*b))
			if finished, isReadMore, err = rf.Read2File(b); err != nil {
				return err
			}
			if finished {
				rf.conn.Write([]byte("CLOSE"))
				util.Debug("close file")
				rf.status = RfStatusInit
				rf.recvOffset = 0
			}
		default:
			return errors.New("undefined recv file status")
		}
		if isReadMore {
			break
		}

	}
	return
}

func recvConn(conn net.Conn, rf *RecvFiles) {
	defer conn.Close()
	var buf [10000]byte
	bytes := make([]byte, 0, 1000)
	defer func() {
		conn.Close()
	}()
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
		bytes = append(bytes, buf[:n]...)
		if err = rf.handleRecvBytes(&bytes); err != nil {
			util.Perrorf("read err:%+v\n", err)
			break
		}
		if true {
			rf.debugRecvBytes(buf[:n], bytes)
		}
	}
}

func (rf *RecvFiles) debugRecvBytes(buf, bytes []byte) {
	n := len(buf)
	rf.recvRealOffset += int64(n)
	rn := rf.recvRealOffset
	util.Debugf("new bytes:%v,real recv offset:%v, unhandled bytes:%d", n, rn, len(bytes))
	// time.Sleep(30 * time.Millisecond)
}

func main() {
	// ??????TCP ????????????
	args := getCliArgs()
	listener, err := net.Listen("tcp", "0.0.0.0:"+args.Port)
	if err != nil {
		fmt.Println("Listen tcp server failed,err:", err)
		os.Exit(1)
	}
	fmt.Println("Listen port:", args.Port)
	fmt.Println("Home:", args.Home)

	for {
		// ??????socket??????
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listen.Accept failed,err:", err)
			continue
		}

		// ??????????????????
		rf := &RecvFiles{
			status: RfStatusInit,
			home:   args.Home,
			force:  args.Force,
			conn:   conn,
		}
		go recvConn(conn, rf)
	}
}
