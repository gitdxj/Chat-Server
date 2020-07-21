package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

type JsonNetAddrConfig struct {
	Addr string
	Port string
}

var JsonConfig *JsonNetAddrConfig


func InitJsonConfig(filename string){
	JsonConfig = new(JsonNetAddrConfig)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = json.Unmarshal(data, JsonConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func GetAddr() string {
	return net.JoinHostPort(JsonConfig.Addr, JsonConfig.Port)
}

