package file

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func IsExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func AbsPath(path string) string {
	var err error
	path, err = exec.LookPath(path)
	if err != nil {
		panic(err)
	}
	path, _ = filepath.Abs(path)

	return filepath.Clean(path)
}

func IsDir(path string) bool {
	if fileInfo, err := os.Stat(path); err != nil {
		return false
	} else {
		return fileInfo.IsDir()
	}
}

func GetParentDir(path string) string {
	return filepath.Dir(path)
}
func GetFilename(path string) string {
	_, filename := filepath.Split(path)
	return filename
}

// Make path force
func MakeFileDir(filePath string) error {
	fileDir := filepath.Dir(filePath)
	return os.MkdirAll(fileDir, os.ModePerm)
}

func CreateFile(filePath string, force bool) (fp *os.File, err error) {
	if !force && IsExists(filePath) {
		return nil, errors.New("file existed:" + filePath)
	}
	if err := MakeFileDir(filePath); err != nil {
		return nil, err
	}
	// fp, err = os.Create(file)
	fp, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	// defer f.Close()
	return

}

func WriteFile(file string, bytes []byte) error {
	err := ioutil.WriteFile(file, bytes, 0644)
	return err
}

func ReadFile(file string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(file)
	return bytes, err
}

func ReadSeek(filePath string, offset int64, n int64) ([]byte, error) {
	fp, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	if _, err := fp.Seek(offset, 0); err != nil {
		return nil, err
	}
	buf := make([]byte, n)
	if m, err := io.ReadFull(fp, buf); err != nil {
		return nil, err
	} else {
		return buf[:m], nil
	}
}

func ReadFileSeek(fp *os.File, offset int64, n int64) ([]byte, error) {
	if _, err := fp.Seek(offset, 0); err != nil {
		return nil, err
	}
	buf := make([]byte, n)
	if _, err := io.ReadFull(fp, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

/**
	f, err := os.Create("/tmp/dat2")
    defer f.Close()
*/
func WriteBytes(f *os.File, buf []byte) (err error) {
	n2, err := f.Write(buf)
	if err != nil {
		return err
	}
	if n2 != len(buf) {
		return fmt.Errorf("write bytes too less:%d!=%d", n2, len(buf))
	}
	return err
}

// return string(a)==string(b)
func BytesDiff(a, b []byte) bool {
	al := len(a)
	bl := len(b)
	if al != bl {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

/**
  -1 no diff
  0 diff at begin
  1 diff at 1 byte
    a:=[]byte("abc")
    fmt.Println(BytesDiffn(a,[]byte("abc"))==-1)
    fmt.Println(BytesDiffn(a,[]byte("abcde"))==3)
    fmt.Println(BytesDiffn(a,[]byte("abde"))==2)
    fmt.Println(BytesDiffn(a,[]byte("ab"))==2)
*/
func BytesDiffn(a, b []byte) int {
	if len(a) > len(b) {
		a, b = b, a
	}
	i, v := 0, byte(0)
	for i, v = range a {
		if v != b[i] {
			return i
		}
	}
	if len(a) == len(b) {
		return -1
	} else {
		return i + 1
	}
}
