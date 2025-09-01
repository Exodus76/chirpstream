package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`   // Used for SQL DBs
	Keyspace string `yaml:"keyspace"` // Used for Scylla/Cassandra
}

type Config struct {
	Databases struct {
		Users      DBConfig `yaml:"users"`
		Users_Test DBConfig `yaml:"users_test"`
	} `yaml:"databases"`
}

func Init(configPath string) (config Config, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}

func (db *DBConfig) DBConnstring() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		db.Driver, db.User, db.Password, db.Host, int(db.Port), db.DBName,
	)
}
