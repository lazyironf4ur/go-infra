package conf

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestOsArgs(t *testing.T) {
	for idx, v := range os.Args {
		fmt.Printf("%d:%s ", idx, v)
	}
	fmt.Println()
}

func TestConf(t *testing.T) {
	fmt.Println(GlobalConfig)
}

func TestString(t *testing.T) {
	s := new(string)
	*s = "abc"
	
	v := reflect.ValueOf(s).Elem()
	v.SetString("daasdasd")
	fmt.Println(v.String())
}

func TestDirect(t *testing.T) {

	type a struct {
		Name string
	}
	s := &a{}
	s.Name = "abc"
	
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
 	field := v.FieldByName("Name")
	field.SetString("daasdasd")
	fmt.Println(s.Name)
}

func TestMap(t *testing.T) {

	m := map[string]interface{}{
		"a": "ABC",
		"b": "ABCD",
	}
	
	v := reflect.ValueOf(m)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() == reflect.Map {
		iter := v.MapRange()
		for iter.Next() {
			k := iter.Key()
			//v2 := iter.Value()
			v.SetMapIndex(k, reflect.ValueOf("hahaha"))
		}
	}

	fmt.Println(m)
	
}