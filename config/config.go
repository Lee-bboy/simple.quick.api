package config

import (
	"log"

	"github.com/kylelemons/go-gypsy/yaml"
)

var Config *yaml.File

func init() {
	var err error
	Config, err = yaml.ReadFile("conf.yaml")
	if err != nil {
		log.Fatal(err)
	}
}
