package configs

import "github.com/spf13/viper"

type Conf struct {
	PORT               string `mapstructure:"PORT"`
	RedisHost          string `mapstructure:"REDIS_HOST"`
	RedisPort          string `mapstructure:"REDIS_PORT"`
	RedisPassword      string `mapstructure:"REDIS_PASSWORD"`
	RedisDatabaseIndex int    `mapstructure:"REDIS_DATABASE_INDEX"`
}

func LoadConfig(path string) (cfg *Conf, err error) {
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return cfg, err
}
