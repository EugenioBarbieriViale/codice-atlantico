package database

import (
	"fmt"
	"strings"
	"os"
	"reflect"
)

type Config struct {
	User string
	Password string
	Dbname string
	Host string
	Port int
	Sslmode string
}

func DefaultConfig() Config {
	return Config{
		User     : os.Getenv("DB_USER"),
		Password : os.Getenv("DB_PASS"),
		Dbname   : "maindb",
		Host     : "localhost",
		Port     : 5432,
		Sslmode  : "disable",
	}
}

func (c Config) ToString() string {
	v := reflect.ValueOf(c)
	t := reflect.TypeOf(c)

	config_str := ""
	for i := 0; i < v.NumField(); i++ {
		if name_str := t.Field(i).Name; i != v.NumField() {
			config_str += strings.ToLower(name_str) + "="
			config_str += fmt.Sprintf("%v", v.Field(i).Interface()) + " "
		}
	}

	return config_str
}
