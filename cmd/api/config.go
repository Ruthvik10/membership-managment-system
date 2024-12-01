package main

import "github.com/spf13/viper"

type config struct {
	DBURL   string `mapstructure:"DB_URL"`
	APIAddr string `mapstructure:"API_ADDR"`
}

func newConfig(path string) (*config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// Enable reading from environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
