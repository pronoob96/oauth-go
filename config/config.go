package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// YamlConfig represents the YAML configuration read from a yml file
type YamlConfig struct {
	Server     Server     `yaml:"server"`
	Datasource Datasource `yaml:"datasource"`
	Redis      Redis      `yaml:"redis"`
}

type Server struct {
	Address string `yaml:"address"`
	Port    int64  `yaml:"port"`
}

type Datasource struct {
	ConnectionURL string `yaml:"connectionurl"`
	Host          string `yaml:"host"`
	Port          int64  `yaml:"port"`
	Databasename  string `yaml:"dbname"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	Authsource    string `yaml:"authsource"`
}

type Redis struct {
	ConnectionURL string `yaml:"connectionurl"`
}

var yamlConfig *YamlConfig

func GetYamlServiceConfig() *YamlConfig {
	config := &YamlConfig{}
	yamlFile, err := ioutil.ReadFile("oauth_local.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return config
}

func GetConfig() *YamlConfig {
	return yamlConfig
}
