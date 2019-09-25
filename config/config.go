package config

import "fmt"

type Config struct {
	Application Application
	Database    Database
}

type Application struct {
	Host string
	Port string
}

type Database struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func create(app Application, db Database) *Config {
	return &Config{Application: app, Database: db}
}

func NewDefault() *Config {
	return create(Application{Host: "0.0.0.0", Port: "3001"}, Database{Host: "0.0.0.0", Port: "5432", Name: "antr", User: "postgres", Password: ""})
}

func (c *Application) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (c *Database) GetAddr() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Name)
}
