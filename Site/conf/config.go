package conf

import (
	"io/ioutil"
	"encoding/json"
)

var (
	Cfg Config
)

type Config struct {
	ConnectionString string `json:"connection_string"`
	ListenPort int `json:"listen_port"`
	RFC2822Format string `json:"rfc_2822_format"`
}

func LoadConfig(path string) (*Config, error){
	if bytes, err := ioutil.ReadFile(path); err != nil {
		return nil, err
	} else {
		conf := Config{}
		if  err := json.Unmarshal(bytes, &conf); err != nil {
			return nil, err
		} else {
			return &conf, err
		}
	}
}