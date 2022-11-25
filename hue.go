package main

import (
	"sync"

	"github.com/amimof/huego"
	"github.com/rs/zerolog/log"
)

type HuegoAdapterDependencies struct {
	Host           string
	User           string
	playerLightMap map[string][]int
}

type HuegoAdapter struct {
	*huego.Bridge
	playerLightMap map[string][]int
	lastState      map[string]map[int]bool
	stateLock      sync.Mutex
}

func NewHuegoAdapter(d *HuegoAdapterDependencies) *HuegoAdapter {
	bridge := huego.New(d.Host, d.User)
	return &HuegoAdapter{
		Bridge:         bridge.Login(d.User),
		lastState:      make(map[string]map[int]bool),
		playerLightMap: d.playerLightMap,
	}
}

// StartSession turns off all the lights configured for that player, and stores all state changes.
func (h *HuegoAdapter) StartSession(playerID string) {
	h.stateLock.Lock()
	defer h.stateLock.Unlock()
	h.lastState[playerID] = h.switchLights(false, h.playerLightMap[playerID]...)
}

// StopSession restores the previous state for all lamps affected by the player uuid session.
func (h *HuegoAdapter) StopSession(playerID string) {
	h.stateLock.Lock()
	defer h.stateLock.Unlock()
	state, ok := h.lastState[playerID]
	if !ok {
		return
	}
	h.setStates(state)
}

// getLights returns all lights by id, returns all lights if no id is specified.
func (h *HuegoAdapter) getLights(lightIDs ...int) (lights []huego.Light) {
	var err error
	lights = make([]huego.Light, 0, len(lightIDs))
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
	return
}

// switchLights change the state ( true = on, false = off ) for all lightIDs.
func (h *HuegoAdapter) switchLights(state bool, lightsIDs ...int) (oldState map[int]bool) {
	var err error
	oldState = make(map[int]bool)
	lights := h.getLights(lightsIDs...)
	for _, light := range lights {
		isOn := light.IsOn()
		// no need to do anything, so we don't register a state change.
		if isOn == state {
			continue
		}
		oldState[light.ID] = isOn
		if state {
			err = light.On()
		} else {
			err = light.Off()
		}
		if err != nil {
			log.Error().Err(err).Msgf("failed to switch lights %d to %v", state)
		}
		log.Debug().Msgf("switching lamp %d to %v", light.ID, state)
	}
	return
}

// setStates receives a map of lightID and state ( on = true, off = false )
// and set it for each light.
func (h *HuegoAdapter) setStates(states map[int]bool) {
	for lightID, state := range states {
		h.switchLights(state, lightID)
	}
	log.Debug().Interface("states", states).Msg("restoring states")
}

func (h *HuegoAdapter) getStates() (states map[int]bool) {
	states = make(map[int]bool)
	for _, light := range h.getLights() {
		states[light.ID] = light.IsOn()
	}
	return
}
