package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct{
	UserName string `json:username`
	PassWord string `json:password`
}



func NewConfig () *Config {
	return &Config{}
}

func (self *Config) Load(filename string) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	datajson := []byte(data)
	err = json.Unmarshal(datajson, self)
	if err != nil {
		return
	}
}
