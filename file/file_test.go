package file

import (
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
