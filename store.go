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
