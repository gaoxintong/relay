package convcode

import (
	"fmt"
	"testing"
)

func TestHex2Dec(t *testing.T) {
	s := "ff"
	dec, err := Hex2Dec(s)
	if err != nil {
		fmt.Println("err: ", err)
	}
	fmt.Println("dec: ", dec)
}

func TestDec2Hex(t *testing.T) {
	i := 123
	s := Dec2Hex(i)
	fmt.Println("hex: ", s)
}
