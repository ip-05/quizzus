package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Google  *GoogleConfig
	Server  *ServerConfig
	Secrets *SecretConfig
}

type ServerConfig struct {
	Host string
	Port int64
}

type GoogleConfig struct {
	ClientId     string
	ClientSecret string
}

type SecretConfig struct {
	Jwt string
}

var config *Config

func GetConfig() *Config {
	if config != nil {
		return config
	}

	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config: %s", err.Error())
	}

	serverConfig := ServerConfig{
		Host: viper.Get("server.host").(string),
		Port: viper.Get("server.port").(int64),
	}

	googleConfig := GoogleConfig{
		ClientId:     viper.Get("google.client_id").(string),
		ClientSecret: viper.Get("google.client_secret").(string),
	}

	secretConfig := SecretConfig{
		Jwt: viper.Get("secrets.jwt").(string),
	}

	config = &Config{
		Server:  &serverConfig,
		Google:  &googleConfig,
		Secrets: &secretConfig,
	}

	return config
}
