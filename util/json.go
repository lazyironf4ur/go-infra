package util

import "encoding/json"

type JsonUtil struct{}


func (jsonUtil JsonUtil) ConvertObjToString(obj interface{}) (string, error){
	b, err := json.Marshal(obj)
	return string(b), err
}