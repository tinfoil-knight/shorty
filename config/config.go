package config

import (
	"log"

	"github.com/spf13/viper"
)

// Get : Gets config variables from a file named config stored at root
func Get(key string) string {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	value, ok := viper.Get(key).(string)
	if ok != true {
		log.Printf("Error in Configuration of %s key\n", key)
	}
	return value
}
