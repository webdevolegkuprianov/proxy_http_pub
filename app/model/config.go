package model

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

//config yaml struct
type Config struct {
	APIVersion string `yaml:"apiVersion"`
	Spec       struct {
		Ports struct {
			Name string `yaml:"name"`
			Addr string `yaml:"bind_addr"`
		} `yaml:"ports"`
		ProxyAddr struct {
			AddrServerRest           string `yaml:"addr_server_rest"`
			AddrServerRestAutoretail string `yaml:"addr_server_rest_ar"`
		} `yaml:"proxy_addr"`
		WhiteIp []string `yaml:"white_ip,flow"`
	} `yaml:"spec"`
}

//New config
func NewConfig() (*Config, error) {

	var service *Config

	f, err := filepath.Abs("/root/config/proxy.yaml")
	if err != nil {
		return nil, err
	}

	y, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(y, &service); err != nil {
		return nil, err
	}

	return service, nil

}
