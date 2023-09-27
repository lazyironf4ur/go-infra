package http

import (
	"fmt"
	"net"
	"testing"
)


func TestIP(t *testing.T) {
	
	a, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, ipaddr := range a {
		if ip, ok := ipaddr.(*net.IPNet); ok && !ip.IP.IsLoopback(){
			fmt.Println(ip.IP.To4().String())
		}
	}
}

func TestFormatIP(t *testing.T) {
	fmt.Println(formatIP())
}

func TestGenerateTraceId(t *testing.T) {
	fmt.Println(generateTraceId())
}