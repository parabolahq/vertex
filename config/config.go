package config

import (
	"errors"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"log"
	"os"
	"strings"
)

var Config = koanf.New(".")

func LoadConfigs() {
	// Loading default values in configuration
	Config.Load(confmap.Provider(map[string]interface{}{
		"bindaddr": "localhost:7000",
		"amqp": map[string]interface{}{
			"url":   "amqp://guest:guest@localhost:5672/",
			"queue": "vertex",
		},
		"config": "config.yaml",
		"keys": []string{
			"public.jwk",
		},
		"id": "vertex-default",
		"handlers": map[string]interface{}{
			"connect":    []string{},
			"disconnect": []string{},
		},
	}, "."), nil)

	// Loading values from Env variables
	Config.Load(env.Provider("VERTEX", "_", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "VERTEX_")), "_", ".", -1)
	}), nil)

	// Loading values from yaml configuration
	_, e := os.Stat(Config.String("config"))
	if !errors.Is(e, os.ErrNotExist) {
		if err := Config.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
			log.Fatalf("error loading config: %v", err)
		}
	}
	log.Println("Config loaded")
}
