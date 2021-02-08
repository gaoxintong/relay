package utils

import (
	"testing"
)

func TestBinaryToHex(t *testing.T) {
	b, _ := BinaryToHex("01010101")
	t.Log("b:", b)
}
