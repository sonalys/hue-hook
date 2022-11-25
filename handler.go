package main

import (
	"encoding/json"
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
	if err := c.Request.ParseMultipartForm(0); err != nil {
		log.Err(err).Msgf("Error while reading form: %v\n", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var payload Payload
	buf, hasPayload := c.Request.MultipartForm.Value["payload"]
	if hasPayload {
		if err := json.Unmarshal([]byte(buf[0]), &payload); err != nil {
			log.Err(err).Msgf("Error while parsing json: %v\n", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
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
