package main

import (
	"fileparser/internal/parsers"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	var configLocation string
	if len(os.Args) > 0 {
		configLocation = os.Args[1]
	} else if os.Getenv("CONFIG_LOCATION") != "" {
		configLocation = os.Getenv("CONFIG_LOCATION")
	} else {
		configLocation = "config.json"
	}
	log.Info("Starting file parser reading config from ", configLocation)
	config := parsers.FileParserConfig{}
	err := config.LoadConfig(configLocation)

	if err != nil {
		panic(err)
	}
	for _, file := range config.Files {
		log.Info("Parsing and writing file ", file.TargetFile, " from ", file.SourceFile, " of type ", file.TargetType)
		err := file.ParseAndWrite()
		if err != nil {
			panic(err)
		}
	}

}
