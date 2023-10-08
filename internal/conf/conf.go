package conf

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	Port     uint16 `yaml:"port"`
	Enable   bool   `yaml:"enable"`
	TokenKey string `yaml:"token_key"`
}

type MysqlConfig struct {
	Dst string `yaml:"dst"`
}

type MongoConfig struct {
	Addr string `yaml:"addr"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type WechatConfig struct {
	Token       string `yaml:"wechat"`
	CallbackUrl string `yaml:"callback_url"`
}

type ChatGtpConfig struct {
	Token string `json:"token"`
	Addr  string `json:"addr"`
}

type Config struct {
	AppConfig     *AppConfig     `yaml:"application"`
	MysqlConfig   *MysqlConfig   `yaml:"mysql"`
	RedisConfig   *RedisConfig   `yaml:"redis"`
	MongoConfig   *MongoConfig   `yaml:"mongo"`
	WechatConfig  *WechatConfig  `yaml:"wechat"`
	ChatGtpConfig *ChatGtpConfig `yaml:"chatgtp"`
}

func InitConfig(configPath string) (c *Config, err error) {
	c = &Config{}
	url := fmt.Sprintf("%s/application.yaml", configPath)

	if os.Getenv("ENV_CONF") == "dev" {
		url = fmt.Sprintf("%s/application-dev.yaml", configPath)
	}

	b, err := os.ReadFile(url)
	if err != nil {
		return
	}
	if err = yaml.Unmarshal(b, c); err != nil {
		return
	}

	return c, nil
}
