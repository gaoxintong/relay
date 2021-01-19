package utils

import (
	"fmt"
	"testing"
)

func TestBinaryToHex(t *testing.T) {
	b, _ := BinaryToHex("01010101")
	fmt.Println("b:", b)
}
