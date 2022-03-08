package config

import (
	"fmt"
	"os"
)

type Database struct {
	User          string
	Password      string
	Database      string
	Host          string
	Port          string
	EnableSSLMode string
	Dialect       string
}

func newDatabaseConfig() Database {
	return Database{
		User:          os.Getenv("DB_USER"),
		Password:      os.Getenv("DB_PASSWORD"),
		Database:      os.Getenv("DB_NAME"),
		Host:          os.Getenv("DB_HOST"),
		Port:          os.Getenv("DB_PORT"),
		EnableSSLMode: RetrieveSSLMode(),
		Dialect:       RetrieveDialect(),
	}
}

func RetrieveDialect() string {
	return os.Getenv("DB_DIALECT")
}

func RetrieveSSLMode() string {
	if env := os.Getenv("ENV"); env == "local" {
		return "disable"
	}
	return "enable"
}

func CreatePostgresConnectionURL() string {
	conf := newDatabaseConfig()
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Database, conf.EnableSSLMode)
}

func CreateDatabaseConnectionURL() string {
	if RetrieveDialect() == "postgres" {
		return CreatePostgresConnectionURL()
	}
	return ""
}
