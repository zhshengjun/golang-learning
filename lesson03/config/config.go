package main

type Database struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
}

type Server struct {
	Name string  `mapstructure:"name"`
	Port string  `mapstructure:"port"`
	Desc *string `mapstructure:"desc"`
}

type Config struct {
	Database Database `mapstructure:"database"`
	Server   Server   `mapstructure:"server"`
}
