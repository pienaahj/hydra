package hydraConfigurator

import (
	"encoding/xml"
	"fmt"
	"os"
)

// decodeXMLConfig decodes the config to xml
func decodeXMLConfig(v interface{}, filename string) error {
	fmt.Println("Decoding XML")
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	return xml.NewDecoder(file).Decode(v)
}
