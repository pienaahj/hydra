package main

import (
	"fmt"

	h "github.com/pienaahj/hydra/hydraConfigurator"
)

type Confstruct struct {
	TS      string  `name:"testString" xml:"testString" json:"testString"`
	TB      bool    `name:"testBool" xml:"testBool" json:"testBool"`
	TF      float64 `name:"testFloat" xml:"testFloat" json:"testFloat"`
	TestInt int
}

func main() {
	configstruct := new(Confstruct)
	// h.GetConfiguration(h.CUSTOM, configstruct, "configfile.conf")
	h.GetConfiguration(h.XML, configstruct, "configfile.xml")
	// h.GetConfiguration(h.JSON, configstruct, "configfile.json")
	fmt.Println(*configstruct)

	if configstruct.TB {
		fmt.Println("bool is true")
	}

	fmt.Println(float64(4.8 * configstruct.TF))

	fmt.Println(5 * configstruct.TestInt)

	fmt.Println(configstruct.TS)
}
