package common

import (
	"log"
	"reflect"
	"runtime"
)

func Must(obj interface{}) {
	_val := reflect.ValueOf(obj)
	if _val.Kind() == reflect.Pointer {
		if _val.IsNil() {
			internalMust(2, _val.Type().String())		
		}
	}
}

func internalMust(skip int, _type string) {
	// bf := make([]byte, 2048)
	// runtime.Stack(bf, false)
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		log.Fatalf("%s:%d %s cannot be nil\n", file, line, _type)
	}
	// fmt.Println(string(bf))
}
