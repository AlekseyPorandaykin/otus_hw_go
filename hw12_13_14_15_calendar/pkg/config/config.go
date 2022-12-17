package config

import "github.com/spf13/viper"

func CreateConfig(pathToFile, typeFile string, conf interface{}) (interface{}, error) {
	viper.SetConfigFile(pathToFile)
	viper.SetConfigType(typeFile)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	return conf, nil
}
