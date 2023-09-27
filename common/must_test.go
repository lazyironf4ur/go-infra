package common

import "testing"


func TestMust(t *testing.T) {
	var p *string
	Must(p)
}

func TestPrintStack(t *testing.T) {
}