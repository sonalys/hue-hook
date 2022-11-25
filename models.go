package main

type Event string

type d map[string]any

const (
	PauseEvent  Event = "media.pause"
	ResumeEvent Event = "media.resume"
	StartEvent  Event = "media.play"
	StopEvent   Event = "media.stop"
)

type Account struct {
	ID    int    `json:"id"`
	Thumb string `json:"thumb"`
	Title string `json:"title"`
}

type Server struct {
	Title string `json:"title"`
	UUID  string `json:"uuid"`
}

type Player struct {
	Local         bool   `json:"local"`
	PublicAddress string `json:"publicAddress"`
	Title         string `json:"title"`
	UUID          string `json:"uuid"`
}

type Payload struct {
	Event   Event   `json:"event"`
	Account Account `json:"account"`
	Server  Server  `json:"server"`
	Player  Player  `json:"player"`
}

type Config struct {
	User           string           `yaml:"user"`
	BridgeHost     string           `yaml:"bridgeHost"`
	PlayerLightMap map[string][]int `yaml:"players"`
	Port           int              `yaml:"port"`
}
