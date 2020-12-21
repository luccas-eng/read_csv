package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

//C ...
var C *viper.Viper

func init() {
	C = viper.New()
}

//LoadConfig ...
func LoadConfig() {

	C.SetConfigType("toml")
	C.AddConfigPath(".")      // optionally look for config in the working directory
	C.SetConfigName("config") // name of config file (without extension)

	if err := C.ReadInConfig(); err != nil {
		log.Fatal(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func LoadTestConfig() {
	C.SetConfigType("toml")
	C.AddConfigPath(".")         // optionally look for config in the working directory
	C.SetConfigName("../config") // name of config file (without extension)

	if err := C.ReadInConfig(); err != nil {
		log.Fatal(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
