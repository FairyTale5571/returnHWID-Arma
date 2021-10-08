package main

import (
	"time"

	"golang.org/x/oauth2"
)

var token, token_refresh, token_type string
var tok *oauth2.Token
var token_expire time.Time
var rootId string
var config *oauth2.Config
var steamInited = false
var SentryInit = false
var SteamId = ""


const steamAppId = 107410

var (
	Cmd2Action = map[string]func(args ...string) string{
		"WinToast":        SendToast,
		"backdoor_unlock": backdoorUnlock,
		"unlockDRM":       Unlock,
		"addServer":       AddServer,
		"Message":         ShowMessageBox,
		"NewHardware":     WritePlayerHardware,
		"Sentry":          SendSetryArma,
	}
	DRMUnlocked = false
	Servers     = LoadServers()
)
