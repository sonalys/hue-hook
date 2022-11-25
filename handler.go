package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HUE interface {
	StartSession(playerID string)
	StopSession(playerID string)
}

type Handler struct {
	hue HUE
}

func NewHandler(hue HUE) *Handler {
	return &Handler{
		hue: hue,
	}
}

func (h *Handler) plexWebhook(c *gin.Context) {
	payload, err := decodePayload(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	logger := log.With().
		Str("event", string(payload.Event)).
		Str("player_uuid", payload.Player.UUID).
		Logger()
	logger.Debug().Msg("received event")
	switch payload.Event {
	case StartEvent, ResumeEvent:
		h.hue.StartSession(payload.Player.UUID)
	case PauseEvent, StopEvent:
		h.hue.StopSession(payload.Player.UUID)
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
