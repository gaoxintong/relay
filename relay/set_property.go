package relay

import (
	"relay/pkg/utils"
	"strings"
)

// SetState 设置状态
func (r *Relay) SetState(state StateCMDType, nos ...uint8) {
	no, err := getNoHex(nos...)
	if err != nil {
		// log
	}
	cmd := "A0 01 08 2B 00 " + no + " " + string(state) + " 00 00 00 00 00 A7"
	r.sendCommand(cmd)
}

// 获取选择编号 16 进制字符串
func getNoHex(nos ...uint8) (string, error) {
	stateArr := strings.Split("00000000", "")
	for _, no := range nos {
		stateArr[7-(no-1)] = "1"
	}
	stateStr := strings.Join(stateArr, "")
	return utils.BinaryToHex(stateStr)
}
