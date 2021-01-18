package server

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strconv"
	"testing"
)

func TestT(t *testing.T) {
	fmt.Println(parseBinToHex("00000001"))
}

func parseBinToHex(s string) (string, error) {
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

func TestCmd(t *testing.T) {
	addr := "192.168.2.77:2701"
	conn, err := net.Dial("tcp", addr)
	// defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	on1, _ := hex.DecodeString("a001082c1100000000000000a7")  // 第一路闭合
	off1, _ := hex.DecodeString("a001082c1000000000000000a7") // 第一路断开
	ops := map[string][][]byte{
		"on": {
			on1,
		},
		"off": {
			off1,
		},
	}
	conn.Write(ops["off"][0])
	b := make([]byte, 128)
	_, err = conn.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("b", b)
	for {
	}
}
