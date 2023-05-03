package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

type (
	Config struct {
		Redis RedisConfig
	}

	RedisConfig struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Password string `mapstructure:"password"`
	}
)

func (cfg *Config) Init(path string) error {
	if err := cfg.InitConfigFile(path); err != nil {
		log.Error(err)
		return err
	}

	if err := parseEnv(); err != nil {
		log.Error(err)
		return err
	}
	cfg.InitFlags()
	cfg.setFromEnv()
	return nil
}

func (cfg *Config) InitConfigFile(path string) error {

	if err := parseConfigFile(path); err != nil {
		return err
	}

	if err := unmarshal(cfg); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) InitFlags() {
	pflag.String("host", "redis", "host redis")
	pflag.String("port", "1234", "port redis")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

func (cfg *Config) setFromEnv() {
	cfg.Redis.Host = viper.GetString("host")
	cfg.Redis.Port = viper.GetString("port")
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("redis", &cfg.Redis); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func parseConfigFile(filepath string) error {
	path := strings.Split(filepath, "/")

	viper.AddConfigPath(path[0])
	viper.SetConfigName(path[1])

	return viper.ReadInConfig()
}

func parseEnv() error {
	if err := parseRedisFromEnv(); err != nil {
		return err
	}
	return nil
}

func parseRedisFromEnv() error {
	err := error(nil)
	viper.SetEnvPrefix("redis")
	if err = viper.BindEnv("host"); err != nil {
		log.Infoln("empty redis host env config")
		log.Error(err)
	}
	if err = viper.BindEnv("port"); err != nil {
		log.Infoln("empty redis port env config")
		log.Error(err)
	}
	if err = viper.BindEnv("password"); err != nil {
		log.Infoln("empty redis password env config")
		log.Error(err)
	}
	return err
}
