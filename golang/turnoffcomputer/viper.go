package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func SetupViper() {

	viper.SetConfigName("application")
	viper.AddConfigPath("resources")
	configReadErr := viper.ReadInConfig()
	if configReadErr != nil {
		fmt.Println("could not read application settings", configReadErr)
		panic(configReadErr)
	}
}
