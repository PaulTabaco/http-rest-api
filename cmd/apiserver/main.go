package main

import (
	"flag"
	"log"
	"paulTabaco/http-rest-api/internal/app/apiserver"

	"github.com/BurntSushi/toml"
)

// Add posibitity add get config path when run our binary. Start
var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

// add posibitity add get config path when run our binary. End

func main() {
	flag.Parse() // commandline flags parse

	config := apiserver.NewConfig()

	_, err := toml.DecodeFile(configPath, config) // Read config file and set vars
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
