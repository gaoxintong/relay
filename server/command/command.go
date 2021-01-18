package command

import (
	"encoding/hex"
	"relay/server/conn"
	"strings"
)

type Command string

const ON, OFF Command = "ON", "OFF" // 控制命令

func SendCommand(no string, cmdStr string) error {
	cmd := Command(cmdStr)
	if cmd == ON {
		cmd = "01"
	} else if cmd == OFF {
		cmd = "02"
	}
	// 获取八路状态

	// 更改单路状态

	cmdStr = "A0 01 08 2B 00 " + no + " " + string(cmd) + " 00 00 00 00 00 A7"
	cmdStr = strings.ReplaceAll(cmdStr, " ", "")
	cmdByte, err := hex.DecodeString(cmdStr)
	if err != nil {
		return err
	}
	conn.Conn.Write(cmdByte)
	return nil
}
