package main

import (
	"log"
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	Host     string   `toml:"host"`
	Port     string   `toml:"port"`
	User     string   `toml:"user"`
	Password string   `toml:"password"`
	DB       string   `toml:"db"`
	JsonCfg  []string `toml:"jsoncfg"`
	Retent   bool     `toml:"retent"`
}

func NewConfig() *Config {
	c := &Config{}
	return c
}
func ParseConfig(path string) (cfg *Config, err error) {
	if path == "" {
		log.Fatalln("no configuration provided, using default settings")
	}
	log.Printf("Using configuration at: %s\n", path)
	config := NewConfig()
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()
	return config, toml.NewDecoder(f).Decode(&config)
}
