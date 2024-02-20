package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

type (
	Config struct {
		HTTP   HTTP
		Rabbit UrlRabbit
		UpBit  UpBit
	}

	HTTP struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderMegabytes"`
	}

	UpBit struct {
		AccessKey string
		SecretKey string
		WsURL     string
	}

	UrlRabbit struct {
		Username     string
		Password     string
		Host         string
		Port         string
		ErlangCookie string
	}
)

//var CFG *Config

// InitConfig initializes the configuration for the application.
func InitConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading .env file: %v", err)
	}

	if err := parseConfigFile("./", os.Getenv("APP_ENV")); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	setFromEnv(&cfg)

	return &cfg, nil
}

//func unmarshal(cfg *Config) error {
//	return viper.UnmarshalKey("http", &cfg.HTTP)
//}

func setFromEnv(cfg *Config) {
	cfg.Rabbit.Username = os.Getenv("RABBIT_USERNAME")
	cfg.Rabbit.Password = os.Getenv("RABBIT_PASSWORD")
	cfg.Rabbit.Host = os.Getenv("RABBIT_HOST")
	cfg.Rabbit.Port = os.Getenv("RABBIT_PORT")
	cfg.Rabbit.ErlangCookie = os.Getenv("RABBIT_COOKIE")

	cfg.UpBit.WsURL = os.Getenv("UPBIT_URL")
	cfg.UpBit.AccessKey = os.Getenv("UPBIT_ACCESS")
	cfg.UpBit.SecretKey = os.Getenv("UPBIT_SECRET")
}

func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.SetConfigName(env)

	return viper.MergeInConfig()
}
