package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Bitflyer struct {
		APIKey    string `yaml:"api_key"`
		APISecret string `yaml:"api_secret"`
	} `yaml:"bitflyer"`
	Trade struct {
		InvestMoney int `yaml:"invest_money"`
		CutLoss     int `yaml:"cut_loss"`
	} `yaml:"trade"`
}

func GetConfig() Config {
	// 读取配置文件
	data, err := os.ReadFile("../config.yaml")
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	// 解析配置文件
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}

	return config
}
