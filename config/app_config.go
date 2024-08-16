package config

import "gopkg.in/yaml.v3"

type ClientConfig struct {
	Name              string `yaml:"name"`
	ForeignServer     string `yaml:"foreign_server"`
	LocalServer       string `yaml:"local_server"`
	ClientCertificate string `yaml:"client_certificate"`
	Sni               string `yaml:"sni"`
}

type ServerConfig struct {
	Name           string `yaml:"name"`
	ListenOn       string `yaml:"listen_on"`
	ServerOn       string `yaml:"server_on"`
	PublicKeyPath  string `yaml:"public_key_path"`
	PrivateKeyPath string `yaml:"private_key_path"`
}

type Config struct {
	Clients []ClientConfig `yaml:"clients"`
	Servers []ServerConfig `yaml:"servers"`
}

func (cfg *Config) SetConfig(ymlFileContent []byte) error {

	return yaml.Unmarshal(ymlFileContent, cfg)
}
