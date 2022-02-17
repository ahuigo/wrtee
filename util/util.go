package util

import (
	"fmt"
	"os"
)

func Perrorf(format string, ds ...interface{}) {
	fmt.Printf(format+"\n", ds...)
}

func Perror(ds ...interface{}) {
	for _, d := range ds {
		fmt.Printf("%v ", d)
	}
	fmt.Printf("\n")
}

func Fatalf(format string, ds ...interface{}) {
	Perrorf(format, ds...)
	os.Exit(1)
}
func Fatal(ds ...interface{}) {
	for _, d := range ds {
		fmt.Printf("%v ", d)
	}
	fmt.Printf("\n")
	os.Exit(1)
}

func Debugf(format string, ds ...interface{}) {
	fmt.Printf(format+"\n", ds...)
}
func Debug(ds ...interface{}) {
	for _, d := range ds {
		fmt.Printf("%v ", d)
	}
	fmt.Printf("\n")
}
