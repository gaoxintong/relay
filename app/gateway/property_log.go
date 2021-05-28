package gateway

import (
	relay2 "relay/app/relay"
	"time"
)

func PropertyLog(devices Devices) relay2.Middleware {
	return func(relay *relay2.Relay, data relay2.Data) relay2.Data {
		if device, ok := devices[relay.SubDeviceID]; ok {
			device.LRUCache.Add(time.Now().Format("2006-05-04 15:02:01")+string(data.PropertyType), &DataRecord{Data: data, Time: time.Now().Format("2006-05-04 15:02:01")})
		}
		return data
	}
}
