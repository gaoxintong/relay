package relay

import (
	"time"

	"github.com/pkg/errors"
)

// WriteFn 向 conn 写入的方法集合
type WriteFn struct {
	d  time.Duration
	fn func()
}

// WriteLoop 开启N个协程，向连接循环发送命令
func (r *Relay) WriteLoop(wfs []WriteFn) error {
	if r.Conn == nil {
		return errors.Wrap(errors.New("not connected"), "write data failed")
	}
	for _, wf := range wfs {
		go func(wf WriteFn) {
			for {
				select {
				case <-r.closed:
					return
				default:
					wf.fn()
					time.Sleep(wf.d)
				}
			}
		}(wf)
	}
	return nil
}
