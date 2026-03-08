package config

import (
	"encoding/json"
	"os"
	"rudolphmax/vbbmon/internal/api"
)

type Config struct {
  Api api.ApiParams;
  DepartureFetchInterval int;
  MessageFetchInterval int;
}

func Read(path string) (Config, error) {
  dat, err := os.ReadFile(path)

  var config Config

  err = json.Unmarshal([]byte(dat), &config);
 	if err != nil {
    return Config{}, err
 	}

  return config, nil
}
