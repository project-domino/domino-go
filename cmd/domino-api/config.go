package main

import (
	"log"
	"os"

	"github.com/project-domino/domino-go/config"
)

// Config is the configuration for the server.
var Config ConfigType

// ConfigType is the type of the configuration for the server.
type ConfigType struct {
	Database config.Database `toml:"database"`
	HTTP     config.HTTP     `toml:"http"`
}

func init() {
	// Create default config object.
	Config = ConfigType{
		Database: config.DefaultDatabase,
		HTTP:     config.DefaultHTTP,
	}

	// Read config or die.
	if err := config.LoadConfig(&Config, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
