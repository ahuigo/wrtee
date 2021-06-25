package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

//client
func main() {
	var err error
	hostPortPtr := flag.String("h", "127.0.0.1:4600", "write host")
	flag.Parse()

	tails := flag.Args()
	hostPortSlice := strings.Split(*hostPortPtr, ":")
	host := hostPortSlice[0]
	port := "4600"
	if len(hostPortSlice) >= 2 {
		port = hostPortSlice[1]
	}
	fileReader := os.Stdin
	if len(tails) > 0 {
		if tails[0] != "-" {
			if fileReader, err = os.Open(tails[0]); err != nil {
				perror("Invalid path: ", tails[0], err)
				// os.Exit(1)
			}
		}
	}

	println(host, port, fileReader)
	// reader to tcp writer
	for {
		buf := make([]byte, 1000)
		n, err := fileReader.Read(buf)
		if err != nil {
			perror(err)
		}
		if n > 0 {
			println(string(buf))
			perror(err)
			break
			// fileWriter.Write(buf)
		} else {
			break
		}
	}

}

func perror(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(1)
}
