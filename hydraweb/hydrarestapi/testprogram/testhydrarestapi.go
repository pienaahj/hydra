package main

import (
	"log"

	hydratestapi "github.com/pienaahj/hydra/hydraweb/hydrarestapi"
)

func main() {
	err := hydratestapi.RunAPI()
	if err != nil {
		log.Fatal(err)
	}
}
