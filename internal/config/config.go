package config

import (
	"encoding/json"
	"os"
	"rudolphmax/vbbmon/internal/api"
	"rudolphmax/vbbmon/internal/display"
)

// Config represents the app configuration parameters.
type Config struct {
  Api api.ApiParams;
  DepartureFetchInterval int;
  MessageFetchInterval int;
  Display display.DisplayConfig;
}

// Read reads the config from the given path and returns a Config struct or an error if the file could not be read or parsed.
// Returns Config and nil, iff the file was read and parsed successfully. Otherwise, returns an empty Config and an error.
func Read(path string) (Config, error) {
  dat, err := os.ReadFile(path)

  var config Config

  err = json.Unmarshal([]byte(dat), &config);
 	if err != nil {
    return Config{}, err
 	}

  return config, nil
}
