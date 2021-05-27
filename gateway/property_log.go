package gateway

import (
	"relay/relay"
	"time"
)

func PropertyLog(devices Devices) relay.Middleware {
	return func(relay *relay.Relay, data relay.Data) relay.Data {
		if device, ok := devices[relay.SubDeviceID]; ok {
			device.LRUCache.Add(time.Now().Format("2006-05-04 15:02:01")+string(data.PropertyType), &DataRecord{Data: data, Time: time.Now().Format("2006-05-04 15:02:01")})
		}
		return data
	}
}
