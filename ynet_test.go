package ynet

import (
	"fmt"
	"testing"
)

func TestSum(t *testing.T) {
	var a = 3
	var b = 4
	res := sum(a, b)
	fmt.Printf("%d 与%d之和:为%d", a, b, res)
	if res != 7 {
		t.Error("error")
	}
}
