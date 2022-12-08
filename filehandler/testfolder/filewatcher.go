package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	WatchFile("testfile.txt")
}

func WatchFile(fname string) {
	filestat1, err := os.Stat(fname)
	PrintFatalError(err)
	for {
		time.Sleep(1 * time.Second)
		filestat2, err := os.Stat(fname)
		PrintFatalError(err)
		if filestat1.ModTime() != filestat2.ModTime() {
			fmt.Println("File was modified at", filestat2.ModTime())
			filestat1, err = os.Stat(fname)
			PrintFatalError(err)
		}
	}
}
func PrintFatalError(err error) {
	if err != nil {
		log.Fatal("Error happened while processing file: ", err)
	}
}
