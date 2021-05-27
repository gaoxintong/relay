package main

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"relay/gateway"
)

func main() {
	TCPAddress := config.TCPAddress
	IOTHubAddress := config.IOTHubAddress
	productKey := config.ProductKey
	deviceName := config.DeviceName
	version := config.Version
	g := gateway.New(TCPAddress, IOTHubAddress, productKey, deviceName, version)
	if err := g.Run(); err != nil {
		panic(err)
	}
}

type Config struct {
	TCPAddress    string `yaml:"TCPAddress"`
	IOTHubAddress string `yaml:"IOTHubAddress"`
	ProductKey    string `yaml:"productKey"`
	DeviceName    string `yaml:"deviceName"`
	Version       string `yaml:"version"`
}

var config = &Config{}

func init() {
	configPath := flag.String("c", "./config/config.yaml", "set your config path")
	yamlFile, err := ioutil.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		panic(err)
	}
}
