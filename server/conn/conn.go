package conn

import (
	"net"

	"github.com/pkg/errors"
)

var Conn net.Conn

// CheckConnState 检查连接状态
func CheckConnState() error {
	if Conn == nil {
		return errors.New("tcp is not connected")
	}
	return nil
}

// ReadData 从连接中读取数据
func ReadData(len int) ([]byte, error) {
	if err := CheckConnState(); err != nil {
		return nil, errors.Wrap(err, "read data failed")
	}
	data := make([]byte, len)
	if _, err := Conn.Read(data); err != nil {
		return nil, errors.Wrap(err, "read data failed")
	}
	return data, nil
}
