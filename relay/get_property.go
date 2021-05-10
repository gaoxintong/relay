package relay

// Data 发送的数据
type Data struct {
	PropertyType PropertyType
	Data         interface{}
}

// 获取数据
func (r *Relay) getData(data Data) interface{} {
	for _, mw := range r.middlewares {
		data = mw(r, data)
	}
	return data.Data
}

// GetOutputState 获取输出状态
func (r *Relay) GetOutputState() Property {
	return r.getData(Data{PropertyType: OUTPUTSTATE, Data: r.outputState})
}

// GetInputState 获取 8 路输入状态
func (r *Relay) GetInputState() Property {
	return r.getData(Data{PropertyType: INPUTSTATE, Data: r.inputState})
}

// GetTH 获取温湿度
func (r *Relay) GetTH() Property {
	return r.getData(Data{PropertyType: TH, Data: r.th})
}
