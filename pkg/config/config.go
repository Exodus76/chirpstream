package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type ServiceConfig struct {
	Addr string `mapstructure:"addr"`
}

type DBConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`   // Used for SQL DBs
	Keyspace string `mapstructure:"keyspace"` // Used for Scylla/Cassandra
}

type scylladbConfig struct {
	ChirpKeyspace    string   `mapstructure:"chirp_keyspace"`
	RelationKeyspace string   `mapstructure:"relationship_keyspace"`
	Host             []string `mapstructure:"host"`
	Username         string   `mapstructure:"username"`
	Password         string   `mapstructure:"password"`
}

type Config struct {
	Databases struct {
		Users      DBConfig `mapstructure:"users"`
		Users_Test DBConfig `mapstructure:"users_test"`
		Chirps     DBConfig `mapstructure:"chirps"`
	} `mapstructure:"databases"`
	Scylladb      scylladbConfig `mapstructure:"scylladb"`
	API_Gateway   ServiceConfig  `mapstructure:"api_gateway"`
	User_Service  ServiceConfig  `mapstructure:"user_service"`
	Chirp_Service ServiceConfig  `mapstructure:"chirp_service"`
	Rel_Service   ServiceConfig  `mapstructure:"rel_service"`
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
