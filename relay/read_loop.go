package relay

import (
	"fmt"
	"relay/pkg/utils"
	"strings"

	"github.com/pkg/errors"
)

// ReadLoop 开启一个协程，从连接中循环读取数据
func (r *Relay) ReadLoop(byteOrderLen int) error {
	if r.Conn == nil {
		return errors.Wrap(errors.New("not connected"), "read data failed")
	}
	go func() {
		for {
			data := make([]byte, byteOrderLen)
			if _, err := r.Conn.Read(data); err != nil {
				fmt.Println("读取失败", err)
				// TODO log
				continue
			}
			cmd, err := cmdToStringLower(data[3])
			if err != nil {
				// TODO log
				continue
			}
			r.dispatchSaveTask(cmd, data)
		}
	}()
	return nil
}

func cmdToStringLower(b byte) (string, error) {
	cmd, err := utils.ByteToHex(b)
	if err != nil {
		return "", err
	}
	return strings.ToLower(cmd), nil
}

func (r *Relay) dispatchSaveTask(cmd string, data []byte) {
	switch cmd {
	case "aa":
		r.SaveOutputState(data)
	case "1b":
		r.SaveInputState(data)
	case "2a":
		r.SaveTH(data)
	}
}
