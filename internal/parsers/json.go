package parsers

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type JSONParser struct {
}

func (j *JSONParser) ParseAndReplace(f FileReplaceOperation) error {
	data := make(map[string]interface{})
	file, err := os.Open(f.SourceFile)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) == 2 {
			key := parts[0]
			valueStr := parts[1]

			if valueInt, err := strconv.Atoi(valueStr); err == nil {
				data[key] = valueInt
			} else {
				data[key] = valueStr
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error("Error reading source file ", f.SourceFile, " ", err)
	}

	if err := j.updateJSONFile(f.TargetFile, data); err != nil {
		panic(err)
	}
	return nil
}

func (j *JSONParser) updateJSONFile(filePath string, keysToReplace map[string]interface{}) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var data interface{}
	if err := json.Unmarshal(file, &data); err != nil {
		return err
	}

	j.replaceValues(data, keysToReplace)

	updatedJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, updatedJSON, 0644); err != nil {
		return err
	}

	return nil
}

func (j *JSONParser) replaceValues(data interface{}, keysToReplace map[string]interface{}) {
	switch concreteVal := data.(type) {
	case map[string]interface{}:
		for k, v := range concreteVal {
			if newValue, ok := keysToReplace[k]; ok {
				concreteVal[k] = newValue
			}
			j.replaceValues(v, keysToReplace)
		}
	case []interface{}:
		for i := range concreteVal {
			j.replaceValues(concreteVal[i], keysToReplace)
		}
	}
}
