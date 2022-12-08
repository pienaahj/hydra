package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("cfile.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	// set up the defaults
	r.Comment = '#'
	r.Comma = ';'

	/*
		records, err := r.ReadAll()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(records)
	*/

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if pe, ok := err.(*csv.ParseError); ok {
				fmt.Println("bad column:", pe.Column)
				fmt.Println("bad line:", pe.Line)
				fmt.Println("Error reported:", pe.Err)
				// after printing the error message continue to the end of the file
				if pe.Err == csv.ErrFieldCount {
					continue
				}

			}
			log.Fatal(err)
		}
		fmt.Println("CSV Row:", record)

		i, err := strconv.Atoi(record[1])
		if err != nil {
			log.Fatal()
		}
		fmt.Println(i * 4)
	}
}
