package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

func readConfig(path string) (c Config, err error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return c, err
	}
	yaml.Unmarshal(file, &c)
	return
}

func envOr(key, fallback string) string {
	v, ok := os.LookupEnv(key)
	if ok {
		return v
	}
	return fallback
}
