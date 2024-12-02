package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Secrets map[string]struct {
		ApiKey    string `yaml:"api_key"`
		ApiSecret string `yaml:"api_secret"`
	} `yaml:"secrets"`
	Trade struct {
		InvestMoney int `yaml:"invest_money"`
		CutLoss     int `yaml:"cut_loss"`
		SafeMoney   int `yaml:"save_money"`
	} `yaml:"trade"`
	DryRun map[string]map[string]bool `yaml:"dry_run"`
}

func GetConfig() Config {
	// 从环境变量读取配置文件路径
	configPath, exist := os.LookupEnv("CONFIG_PATH")
	if !exist {
		configPath = "config.yaml" // 默认值
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
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
