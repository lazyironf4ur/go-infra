package util

import (
	"runtime"
	"strconv"
	"strings"
)



func GetGoroutineID() (int, error) {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	field := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	return strconv.Atoi(field)
}