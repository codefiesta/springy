package util

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	shared *Configuration
	once   sync.Once
)

type ServerConfiguration struct {
	Port int
}

type DatabaseConfiguration struct {
	Host  string
	Port int
	Name string
}

// Builds the fully qualified URI
func (c *DatabaseConfiguration) Uri() string {
	uri := strings.Join([]string{"mongodb://", c.Host, ":", strconv.Itoa(c.Port)}, "")
	return uri
}

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
}


// Our singleton instance of the Configuration
func Config() *Configuration {

	once.Do(func() {

		path, _ := os.Getwd()
		fmt.Print(path)
		log.Println("Configuring ...")
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("..")
		viper.AddConfigPath(path)
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}
		err := viper.Unmarshal(&shared)

		if err != nil {
			log.Fatalf("unable to decode into struct, %v", err)
		}
	})
	return shared
}
