package utils

import (
	"bytes"
	"fmt"
	"strconv"
)

// ByteToBinary byte 转二进制
func ByteToBinary(b byte) string {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(fmt.Sprintf("%08b", b))
	return buf.String()
}

// BinaryToHex 二进制转十六进制
func BinaryToHex(s string) (string, error) {
	ui, err := strconv.ParseUint(s, 2, 64)
	if err != nil {
		return "error", err
	}
	ret := fmt.Sprintf("%x", ui)
	if len(ret) == 1 {
		ret = "0" + ret
	}
	return ret, nil
}

// ByteToHex byte 转十六进制
func ByteToHex(b byte) (string, error) {
	return BinaryToHex(ByteToBinary(b))
}
