package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/amimof/huego"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func main() {
	var lights bool
	flag.BoolVar(&lights, "lights", false, "set this flag for getting all lights from current bridge")
	flag.Parse()

	configPath := envOr("CONFIG", "config.yaml")
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	c, err := readConfig(configPath)
	if err != nil {
		log.Info().Msg("failed to read config.yaml")
		log.Info().Msg("discovering new hue bridge devices on your network")
		bridge, err := huego.Discover()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to discover new bridge devices")
		}
		log.Info().Msg("press the hue bridge device button")

		var user string
		for {
			user, err = bridge.CreateUser("hue-hook")
			if err == nil {
				break
			}
			time.Sleep(time.Second)
		}
		log.Info().Msg("generating new config file")
		c = Config{
			User:       user,
			BridgeHost: bridge.Host,
			Port:       3304,
		}
		buf, err := yaml.Marshal(c)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to encode config")
		}
		err = os.WriteFile(configPath, buf, os.ModeDevice)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to write config")
		}
	}
	bridge := NewHuegoAdapter(c.BridgeHost, c.User)

	if lights {
		log.Info().Interface("lights", bridge.GetStates())
		return
	}

	log.Info().Interface("configuredPlayers", c.PlayerLightMap).Send()
	handler := NewHandler(bridge, c.PlayerLightMap)
	handler.Run(fmt.Sprintf(":%d", c.Port))
}
