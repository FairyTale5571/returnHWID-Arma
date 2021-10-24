package main

import (
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
)

var (
	logFile                          *os.File
	logger                           *log.Logger
	tok                              *oauth2.Token
	config                           *oauth2.Config
	token_expire                     time.Time
	tempDir                          string
	rootId                           string
	token, token_refresh, token_type string
	IV_PUBLIC                        string
	IV_PRIVATE                       string
	DRMUnlocked                      = false
	steamInited                      = false
	SentryInit                       = false
	InfiInited                       = false
	SteamId                          = ""
	dataBus                          chan ExtData
)

const (
	InfistarCloud = "https://api.infistar.vision/v1/cloudban/"
	LKAPI         = "http://localhost:3000/api/admin/"
	steamAppId    = 107410
)

var (
	Cmd2Action = map[string]func(args []string) string{
		"backdoor_unlock": backdoorUnlock,
		"unlockDRM":       Unlock,
		"addServer":       AddServer,
		"Message":         ShowMessageBox,
		"NewHardware":     WritePlayerHardware,
		"Sentry":          SendSetryArma,
	}
	Servers = LoadServers()
)

type Vision struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type Ban struct {
	Ban    string    `json:"ban"`
	Reason string    `json:"ban_reason"`
	Time   time.Time `json:"ban_time"`
}

type ExtData struct {
	Func string
	Args []string
}
