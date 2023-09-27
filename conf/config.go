package conf

import (
	"fmt"
	"os"
	"os/user"
	"reflect"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type ConfigOpt func(interface{}) error

var defaultConfigFilePath = "./config.yaml"

var ParseOsEnvHandler ConfigOpt = parseOsEnv

var GlobalConfig map[string]interface{}

var cacheConfig map[string]interface{} = make(map[string]interface{})

func init() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	defaultConfigFilePath = user.HomeDir + "/config/" + "config.yaml"
	var filePath = defaultConfigFilePath
	for _, v := range os.Args {
		if strings.Contains(v, "config") {
			kv := strings.Split(v, "=")
			filePath = kv[1]
		}
	}
	ReadFileConf(filePath, nil)
	GlobalConfig = make(map[string]interface{})
	for k, v := range cacheConfig {
		GlobalConfig[k] = v
	}
}

func ReadFileConf(filepath string, obj interface{}) {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal(bytes, cacheConfig)
	ConfPostHandle(cacheConfig, []ConfigOpt{parseOsEnv})
	if obj != nil {
		yaml.Unmarshal(bytes, obj)
		ConfPostHandle(obj, []ConfigOpt{parseOsEnv})
	}
}

func ConfPostHandle(conf interface{}, opts []ConfigOpt) {
	for _, opt := range opts {
		err := opt(conf)
		if err != nil {
			panic(err)
		}
	}
}

func parseOsEnv(conf interface{}) error {
	r := regexp.MustCompile(`\$\{\w+\}`)
	_value := reflect.ValueOf(conf)
	return deepParse(_value, r)
}

func deepParse(v reflect.Value, reg *regexp.Regexp) error {

	switch v.Kind() {
	case reflect.Pointer:
		return parsePointer(v, reg)
	case reflect.Map:
		return parseMap(v, reg)
	case reflect.Interface:
		return parseInterface(v, reg)
	default:
		return fmt.Errorf("error type %s to parse config", v.Type().String())
	}
}

func parseAndSet(v reflect.Value, reg *regexp.Regexp) {
	_var := reg.FindString(v.String())
	if _var != "" {
		s := os.Getenv(_var[2 : len(_var)-1])
		v.SetString(s)
	}
}

func parsePointer(v reflect.Value, reg *regexp.Regexp) error {
	v = v.Elem()
	fieldsNum := v.NumField()
	for i := 0; i < fieldsNum; i++ {
		fieldValue := v.Field(i)
		//fmt.Printf("parsing %s\n", fieldValue.Type().String())
		if fieldValue.Kind() == reflect.Struct || fieldValue.Kind() == reflect.Pointer {
			err := deepParse(fieldValue, reg)
			if err != nil {
				return err
			}
			continue
		}

		if fieldValue.Kind() == reflect.String {
			parseAndSet(fieldValue, reg)
		} else {
			return fmt.Errorf("parsing config failed, unsupported type of %s", v.Type().String())
		}
	}

	return nil
}

func parseMap(v reflect.Value, reg *regexp.Regexp) error {
	iter := v.MapRange()
	for iter.Next() {
		_k := iter.Key()
		_v := iter.Value()
		if _v.Kind() == reflect.String {
			_var := reg.FindString(_v.String())
			if _var != "" {
				s := os.Getenv(_var[2 : len(_var)-1])
				v.SetMapIndex(_k, reflect.ValueOf(s))
			}
		} else if _v.Kind() == reflect.Interface && _v.Elem().Kind() == reflect.String {
			_var := reg.FindString(_v.Elem().String())
			if _var != "" {
				s := os.Getenv(_var[2 : len(_var)-1])
				v.SetMapIndex(_k, reflect.ValueOf(s))
			}
		} else {
			err := deepParse(_v, reg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func parseInterface(v reflect.Value, reg *regexp.Regexp) error {
	_v := v.Elem()
	if _v.Kind() == reflect.String {
		parseAndSet(_v, reg)
	} else if _v.Kind() == reflect.Map {
		if err := deepParse(_v, reg); err != nil {
			return err
		}

	}

	return nil
}
