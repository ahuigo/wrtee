package file

import (
	"fmt"
	"testing"
)

func TestCreateFile(t *testing.T) {
	file := "tmp/a.txt"
	fp, _ := CreateFile(file, true)

	WriteBytes(fp, []byte("hello"))
	WriteBytes(fp, []byte(" world"))
	fp.Close()

	r, _ := ReadFile(file)

	if string(r) != "hello world" {
		t.Fatal("write failed")
	}

}

func TestReadSeek(t *testing.T) {
	buf, err := ReadSeek("../src/wget", 0, 10)
	fmt.Println("err:", err)
	o := buf
	fmt.Println("output1:", string(o))

	buf, _ = ReadSeek("../src/wget", 100, 10)
	o = append(o, buf...)
	fmt.Println("output2:", string(o))

}
