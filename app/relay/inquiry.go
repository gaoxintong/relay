package relay

// InquiryTH 发送查询温湿度命令
func (r *Relay) InquiryTH() {
	cmd := "A0 01 08 2A 00 00 00 00 00 00 00 00 A7"
	r.sendCommand(cmd)
}

// InquiryInputState 发送查询输入状态命令
func (r *Relay) InquiryInputState() {
	cmd := "A0 01 08 1B 00 00 00 00 00 00 00 00 A7"
	r.sendCommand(cmd)
}
