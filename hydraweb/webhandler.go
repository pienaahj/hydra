package hydraweb

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
)

type testhandler struct {
	r int
}

// Create a constructor
func newHandler() testhandler {
	rInt, err := rand.Int(bytes.NewReader([]byte("testbyte")), big.NewInt(10))
	if err != nil {
		log.Fatal("cannot get INT", err)
	}
	if !rInt.IsInt64() {
		log.Fatal("cannot get INT", err)
	}
	// convert it to int
	r := int(rInt.Int64())
	return testhandler{r: r}
}

func (h testhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprint(w, "Welsome to the Hydra software system from the custom handler")
	case "/testhandle":
		fmt.Fprint(w, "test handle object with random number: ", h.r)
	}
	fmt.Println(r.URL.Query())
}
