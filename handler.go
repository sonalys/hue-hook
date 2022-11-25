package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HUE interface {
	SwitchLights(state bool, lights ...int)
	GetStates(lights ...int) (states map[int]bool)
	RestoreStates(states map[int]bool)
}

type Handler struct {
	hue            HUE
	lastState      map[int]bool
	playerLightMap map[string][]int
}

func NewHandler(hue HUE, playerLightMap map[string][]int) *Handler {
	return &Handler{
		hue:            hue,
		playerLightMap: playerLightMap,
		lastState:      make(map[int]bool),
	}
}

func (h *Handler) plexWebhook(c *gin.Context) {
	payload, err := decodePayload(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	lights := h.playerLightMap[payload.Player.UUID]
	logger := log.With().
		Str("event", string(payload.Event)).
		Str("player_uuid", payload.Player.UUID).
		Interface("lamps", lights).
		Logger()
	logger.Debug().Msg("received event")
	switch payload.Event {
	case StartEvent, ResumeEvent:
		h.lastState = h.hue.GetStates(lights...)
		h.hue.SwitchLights(false, lights...)
	case PauseEvent, StopEvent:
		h.hue.RestoreStates(h.lastState)
	default:
		logger.Error().Msgf("failed to handle event")
	}
	c.Status(http.StatusOK)
}

func (h *Handler) Run(addr string) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.POST("/", h.plexWebhook)
	return r.Run(addr)
}
