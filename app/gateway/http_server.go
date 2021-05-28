package gateway

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"relay/pkg/convcode"
	"strconv"
)

var r *mux.Router

func init() {
	r = mux.NewRouter()
}

func (g *Gateway) startHTTPServer() {
	r.HandleFunc("/devices", g.allDeviceLog)
	r.HandleFunc("/devices/{id}", g.getDeviceLog)
	http.ListenAndServe("0.0.0.0:6060", r)
}

func (g *Gateway) allDeviceLog(rw http.ResponseWriter, req *http.Request) {
	type Device struct {
		DeviceCodeHex   string      `json:"deviceCodeHex"`
		DeviceCodeAscii uint16      `json:"deviceCodeAscii"`
		Address         string      `json:"address"`
		OnlineTime      string      `json:"onlineTime"`
		Log             interface{} `json:"log"`
	}
	type DeviceInfo struct {
		Count uint64    `json:"count"`
		List  []*Device `json:"list"`
	}

	var devices DeviceInfo
	for _, device := range g.Devices {
		log := []interface{}{}
		for _, Node := range device.LRUCache.GetAll() {
			log = append(log, Node.Value)
		}
		devices.List = append(devices.List, &Device{
			DeviceCodeHex:   convcode.Dec2Hex(int(device.SubDeviceID)),
			DeviceCodeAscii: device.SubDeviceID,
			Address:         device.Conn.RemoteAddr().String(),
			OnlineTime:      device.OnlineTime,
			Log:             log,
		})
	}
	devices.Count = uint64(len(devices.List))
	b, err := json.Marshal(devices)
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}
	rw.Write(b)
}

func (g *Gateway) getDeviceLog(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		rw.Write([]byte("id type error" + err.Error()))
		return
	}
	type Device struct {
		DeviceCodeHex   string      `json:"deviceCodeHex"`
		DeviceCodeAscii uint16      `json:"deviceCodeAscii"`
		Address         string      `json:"address"`
		OnlineTime      string      `json:"onlineTime"`
		Log             interface{} `json:"log"`
	}
	var ret *Device
	device, ok := g.Devices[uint16(idUint)]
	if !ok {
		rw.Write([]byte(err.Error()))
		return
	}
	log := []interface{}{}
	for _, Node := range device.LRUCache.GetAll() {
		log = append(log, Node.Value)
	}
	ret = &Device{
		DeviceCodeHex:   convcode.Dec2Hex(int(device.SubDeviceID)),
		DeviceCodeAscii: device.SubDeviceID,
		Address:         device.Conn.RemoteAddr().String(),
		OnlineTime:      device.OnlineTime,
		Log:             log,
	}

	b, err := json.Marshal(ret)
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}
	rw.Write(b)
}
