package main

import (
	"github.com/amimof/huego"
	"github.com/rs/zerolog/log"
)

type HuegoAdapter struct {
	*huego.Bridge
}

func NewHuegoAdapter(host, user string) *HuegoAdapter {
	bridge := huego.New(host, user)
	return &HuegoAdapter{bridge.Login(user)}
}

func (h *HuegoAdapter) SwitchLights(state bool, lights ...int) {
	for _, lightID := range lights {
		light, err := h.GetLight(lightID)
		if err != nil {
			log.Error().Err(err).Msgf("failed to get light %d", lightID)
			continue
		}
		if state {
			err = light.On()
		} else {
			err = light.Off()
		}
		if err != nil {
			log.Error().Err(err).Msgf("failed to switch lights %d to %v", state)
		}
		log.Debug().Msgf("switching lamp %d to %v", lightID, state)
	}
}

func (h *HuegoAdapter) GetStates(lightIDs ...int) (states map[int]bool) {
	var err error
	var lights []huego.Light
	if len(lightIDs) == 0 {
		lights, err = h.GetLights()
		if err != nil {
			log.Error().Err(err).Msg("failed to get lights from bridge")
		}
	} else {
		for _, lightID := range lightIDs {
			light, err := h.GetLight(lightID)
			if err != nil {
				log.Error().Err(err).Msgf("failed to get light %d", lightID)
				continue
			}
			lights = append(lights, *light)
		}
	}
	states = make(map[int]bool, len(lights))
	for _, light := range lights {
		states[light.ID] = light.IsOn()
	}
	log.Debug().Interface("states", states).Msg("light states")
	return
}

func (h *HuegoAdapter) RestoreStates(states map[int]bool) {
	for lightID, state := range states {
		h.SwitchLights(state, lightID)
	}
	log.Debug().Interface("states", states).Msg("restoring states")
}
