package hydraConfigurator

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	CUSTOM uint8 = iota
	JSON
	XML
)

var wrongTypeError error = errors.New("type must be a pointer to a struct")

func GetConfiguration(confType uint8, obj interface{}, filename string) (err error) {
	//  Check if this is a type pointer
	mysRValue := reflect.ValueOf(obj)
	if mysRValue.Kind() != reflect.Ptr || mysRValue.IsNil() {
		return wrongTypeError
	}
	//  get and confirm the struct value
	mysRValue = mysRValue.Elem()
	//  *object => object
	//  reflection value of *object .Elem() => object() (Settable!!)
	if mysRValue.Kind() != reflect.Struct {
		return wrongTypeError
	}

	switch confType {
	case CUSTOM:
		fmt.Println("CUSTOM case selected...")
		err = MarshalCustomConfig(mysRValue, filename)
	case JSON:
		fmt.Println("JSON case selected...")
		err = decodeJSONConfig(obj, filename)
	case XML:
		fmt.Println("XML case selected...")
		err = decodeXMLConfig(obj, filename)
	}
	return err
}
