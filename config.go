package main

import (
	"encoding/json"
	"os"
)

type Config struct {
  Api ApiParams;
  FetchInterval int;
}

func readConfig(path string) (Config, error) {
  dat, err := os.ReadFile(path)

  var config Config

  err = json.Unmarshal([]byte(dat), &config);
 	if err != nil {
    return Config{}, err
 	}

  return config, nil
}
