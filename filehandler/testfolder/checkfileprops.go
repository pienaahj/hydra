package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	GeneratateFileStatusReport("testfile.txt")
}

func GeneratateFileStatusReport(fname string) {
	// Stat returns file info.  It will return an error if there is no file.
	filestats, err := os.Stat(fname)
	PrintFatalError(err)
	fmt.Println("What's the file name?", filestats.Name())
	fmt.Println("Am i a directory?", filestats.IsDir())
	fmt.Println("What are the permissions?", filestats.Mode())
	fmt.Println("What's the file size?", filestats.Size())
	fmt.Println("When was it last modified?", filestats.ModTime())
}
func PrintFatalError(err error) {
	if err != nil {
		log.Fatal("Error happened while processing file: ", err)
	}
}
