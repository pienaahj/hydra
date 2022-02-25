package main

import (
	"fmt"

	"github.com/pienaahj/hydra/hydrachat"
)

func main() {
	exit := make(chan bool)
	fmt.Println("Server starting at port 2100...")
	hydrachat.Run()
	<-exit
}