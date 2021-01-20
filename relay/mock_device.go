package relay

import (
	"net"
	"time"

	"github.com/pkg/errors"
)

// MockDevice 模拟子设备
type MockDevice struct {
	TCPServerAddress string
	IDM              string
	IDN              string
	KeepAlive        time.Duration
	Conn             net.Conn
}

// InitTCPClient 初始化 TCP 连接
func (m *MockDevice) InitTCPClient() {
	var err error
	if m.Conn, err = net.Dial("tcp4", m.TCPServerAddress); err != nil {
		panic(err)
	}
	cmd := "A0 " + m.IDM + " " + m.IDN + " 2A 00 00 00 00 00 00 00 00 A7"
	m.sendCommand(cmd)
}

// AutoPostDeviceInfo 自动发送设备信息
func (m *MockDevice) AutoPostDeviceInfo() {
	go func() {
		for {
			time.Sleep(m.KeepAlive)
			cmd := "A0 " + m.IDM + " " + m.IDN + " 2A 00 00 00 00 00 00 00 00 A7"
			m.sendCommand(cmd)
		}
	}()
}

func (m *MockDevice) sendCommand(cmd string) error {
	cmdByte, err := commandFormatter(cmd)
	if err != nil {
		return errors.Wrap(err, "send command failed")
	}
	if _, err = m.Conn.Write(cmdByte); err != nil {
		return errors.Wrap(err, "send command failed")
	}
	return nil
}
