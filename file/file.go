package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func IsExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
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
	_, file := filepath.Split(path)
	return file
}

// Make path force
func MakeFileDir(filePath string) error {
	fileDir := filepath.Dir(filePath)
	return os.MkdirAll(fileDir, os.ModePerm)
}

func CreateFile(file string) (fp *os.File, err error) {
	if err := MakeFileDir(file); err != nil {
		return nil, err
	}
	fp, err = os.Create(file)
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