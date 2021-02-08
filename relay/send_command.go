package relay

import (
	"encoding/hex"
	"strings"

	"github.com/pkg/errors"
)

// 发送命令
func (r *Relay) sendCommand(cmd string) error {
	cmdByte, err := commandFormatter(cmd)
	if err != nil {
		return errors.Wrap(err, "send command failed")
	}
	if _, err = r.Conn.Write(cmdByte); err != nil {
		return errors.Wrap(err, "send command failed")
	}
	return nil
}

// 格式化命令
func commandFormatter(originalCommand string) ([]byte, error) {
	originalCommand = strings.ReplaceAll(originalCommand, " ", "")
	return hex.DecodeString(originalCommand)
}
