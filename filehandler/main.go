package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var filename string = "test1.txt"
	var filename2 string = "test2.txt"
	var filename3 string = "test3.txt"
	// Open a file read only
	f1, err := os.Open(filename)
	PrintFatalError(err)
	defer f1.Close()

	// Create a new file
	f2, err := os.Create(filename2)
	PrintFatalError(err)
	defer f2.Close()

	// Open a file to read write
	f3, err := os.OpenFile(filename3, os.O_APPEND|os.O_RDWR, 0666)
	PrintFatalError(err)
	defer f3.Close()

	/*
			const (
			// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
			O_RDONLY int = syscall.O_RDONLY // open the file read-only.
			O_WRONLY int = syscall.O_WRONLY // open the file write-only.
			O_RDWR   int = syscall.O_RDWR   // open the file read-write.
			// The remaining values may be or'ed in to control behavior.
			O_APPEND int = syscall.O_APPEND // append data to the file when writing.
			O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
			O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist.
			O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
			O_TRUNC  int = syscall.O_TRUNC  // truncate regular writable file when opened.
		)
		 0666 => Owner: (read & write), Group (read & write), and other (read & write)

		Symbolic	Numeric			English
		notation	notation

		----------	0000	no permissions
		-rwx------	0700	read, write, & execute only for owner
		-rwxrwx---	0770	read, write, & execute for owner and group
		-rwxrwxrwx	0777	read, write, & execute for owner, group and others
		---x--x--x	0111	execute
		--w--w--w-	0222	write
		--wx-wx-wx	0333	write & execute
		-r--r--r--	0444	read
		-r-xr-xr-x	0555	read & execute
		-rw-rw-rw-	0666	read & write
		-rwxr-----	0740	owner can read, write, & execute; group can only read; others have no permissions
	*/

	// rename a file
	err = os.Rename(filename, filename2)
	PrintFatalError(err)

	// move a file
	err = os.Rename("./test1.txt", "./testfolder/text1.txt")
	PrintFatalError(err)

	// copy a file
	CopyFile("test1.txt", "./testfolder/text3.txt")

	// delete a file
	err = os.Remove(filename)
	PrintFatalError(err)

	bytes, err := ioutil.ReadFile(filename2)
	PrintFatalError(err)
	fmt.Println(string(bytes))

	scanner := bufio.NewScanner(f3)
	count := 0
	for scanner.Scan() {
		count++
		fmt.Println("Found line: ", count, scanner.Text())
	}

	// buffered write, efficient store in memory, saves disk I/O
	writebuffer := bufio.NewWriter(f3)
	for i := 1; i <= 5; i++ {
		writebuffer.WriteString(fmt.Sprintln("Added line", i))
	}
	writebuffer.Flush() // writes the buffer to io.Writer

	GeneratateFileStatusReport(filename3)

	filestat1, err := os.Stat(filename3)
	PrintFatalError(err)
	for {

	}

}
func PrintFatalError(err error) {
	if err != nil {
		log.Fatal("Error happened while processing file: ", err)
	}
}

// CopyFile will copy a file from fname1 to fname2
func CopyFile(fname1, fname2 string) {
	fOld, err := os.Open(fname1)
	PrintFatalError(err)
	defer fOld.Close()

	fNew, err := os.Create(fname2)
	PrintFatalError(err)
	defer fNew.Close()

	// copy bytes from source to destination
	_, err = io.Copy(fNew, fOld)
	PrintFatalError(err)

	// flush the contents to disk(hard drive)
	err = fNew.Sync()
	PrintFatalError(err)
}

func GeneratateFileStatusReport(fname string) {
	// Stat returns file info.  It will return an error if there is no file.
	filestats, err := os.Stat(fname)
	PrintFatalError(err)
	fmt.Println("What's the file name?", filestats.Name())
	fmt.Println("Am i a directory?", filestats.IsDir())
	fmt.Println("What are the permissions?", filestats.Mode())
	fmt.Println("What's the file nsize?", filestats.Size())
	fmt.Println("When was it last modified?", filestats.ModTime())
}
