package hydraConfigurator

import (
	"errors"
	"reflect"
	"fmt"
)

const (
	CUSTOM uint8 = iota
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
	}
	return err
}
