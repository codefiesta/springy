package util

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	env  *Environment
	once sync.Once
)

type ServerEnv struct {
	Port int
}

type DatabaseEnv struct {
	Host       string
	Port       int
	Db         string
	Collection string
	Username   string
	Password   string
	ReplicaSet string
}

type Environment struct {
	Server   ServerEnv
	Database DatabaseEnv
}

// Builds the fully qualified host URI
func (e *DatabaseEnv) GetURI() string {
	uri := strings.Join([]string{
		e.Host,
		":",
		strconv.Itoa(e.Port),
	}, "")
	return uri
}

//Our singleton instance of the Environment
func Env() *Environment {

	once.Do(func() {

		log.Println("🌱 [Configuring Springy] 🌱")
		viper.SetConfigFile(".env")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}

		dir, _ := os.Getwd()
		log.Println("🎯", dir)

		db := DatabaseEnv{
			Host:       viper.GetString("MONGO_HOST"),
			Port:       viper.GetInt("MONGO_PORT"),
			Db:         viper.GetString("MONGO_DB"),
			Collection: viper.GetString("MONGO_COLLECTION"),
			Username:   viper.GetString("MONGO_USER"),
			Password:   viper.GetString("MONGO_PASSWORD"),
			ReplicaSet: viper.GetString("MONGO_REPLICA_SET"),
		}

		server := ServerEnv{
			Port: viper.GetInt("SERVER_PORT"),
		}

		env = &Environment{
			Server:   server,
			Database: db,
		}
	})
	return env
}
