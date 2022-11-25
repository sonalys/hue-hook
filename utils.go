package main

import (
	"encoding/json"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

func decodePayload(c *gin.Context) (p Payload, err error) {
	err = c.Request.ParseMultipartForm(c.Request.ContentLength)
	if err != nil {
		return
	}
	buf := c.Request.MultipartForm.Value["payload"][0]
	err = json.Unmarshal([]byte(buf), &p)
	return
}

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
