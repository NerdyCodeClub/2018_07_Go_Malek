package outletconfiguration

import (
	"encoding/json"
	"os"
)

// Configuration holds API config
type Configuration struct {
	Port           int
	VeSyncEndpoint string
}

// LoadConfiguration loads config values
func LoadConfiguration() Configuration {
	var config Configuration
	filename := "config.json"

	if fileExists("config.development.json") {
		filename = "config.development.json"
	}
	file, error := os.Open(filename)
	if error != nil {
		return config
	}

	decoder := json.NewDecoder(file)
	decoder.Decode(&config)

	return config
}

func fileExists(filename string) bool {
	if _, error := os.Stat(filename); error != nil {
		if os.IsNotExist(error) {
			return false
		}
	}
	return true
}
