package parsers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type PropertiesParser struct {
}

func (p *PropertiesParser) ParseAndReplace(f FileReplaceOperation) error {
	// Load key-value pairs from the source file
	keyValues, err := p.loadPropertiesFile(f.SourceFile)
	if err != nil {
		return err
	}

	// Update target properties file
	return p.updatePropertiesFile(f.TargetFile, keyValues)
}

// loadPropertiesFile reads key-value pairs from a properties file, ignoring lines that start with "#"
func (p *PropertiesParser) loadPropertiesFile(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	keyValues := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue // ignore comments
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			keyValues[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return keyValues, nil
}

// updatePropertiesFile replaces keys in the target properties file with values from keyValues
func (p *PropertiesParser) updatePropertiesFile(filePath string, keyValues map[string]string) error {
	// Read the existing content of the target file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")

	// Replace keys in the file content
	for i, line := range lines {
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue // ignore comments and empty lines
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			if newValue, ok := keyValues[key]; ok {
				lines[i] = fmt.Sprintf("%s=%s", key, newValue)
			}
		}
	}

	// Write the modified content back to the file
	updatedContent := strings.Join(lines, "\n")
	return os.WriteFile(filePath, []byte(updatedContent), 0644)
}
