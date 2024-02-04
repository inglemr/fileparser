package parsers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Parser interface {
	ParseAndReplace(f FileReplaceOperation) error
}

type FileParserConfig struct {
	Files []FileReplaceOperation `json:"files"`
}

type FileReplaceOperation struct {
	TargetFile string `json:"target_file"`
	SourceFile string `json:"source_file"`
	TargetType string `json:"target_type"`
}

func (f *FileParserConfig) LoadConfig(configLocation string) error {
	file, err := os.Open(configLocation)
	if err != nil {
		return err
	}
	err = os.Chdir(filepath.Dir(configLocation))
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&f)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileReplaceOperation) ParseAndWrite() error {
	//We need to parse the source file which will have key value pairs to find in the target file and replace the value of
	switch strings.ToUpper(f.TargetType) {
	case "JSON":
		parser := &JSONParser{}
		return parser.ParseAndReplace(*f)
	case "PROPERTIES":
		parser := &PropertiesParser{}
		return parser.ParseAndReplace(*f)
	default:
		log.Warn("Unknown file type ", f.TargetType, " for file ", f.TargetFile, " Skipping...")
	}
	return nil
}
