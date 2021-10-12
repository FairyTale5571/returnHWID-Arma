package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include "extensionCallback.h"
*/
import "C"

type Vision struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func CheckInfiBan() string {
	if !InfiInited {
		return "Cloud ban need init"
	}
	client := &http.Client{}

	resp, err := http.NewRequest("GET", InfistarCloud+GetPlayerUid()+"/baninfo", nil)
	if err != nil {
		SendSentry(err.Error())
	}
	resp.Header.Set("vi-server-ident", IV_PUBLIC)
	resp.Header.Set("vi-server-secret", IV_PRIVATE)

	req, err := client.Do(resp)
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))

	t := Vision{}
	err = json.Unmarshal(body, &t)
	if err != nil {
		SendSentry(err.Error())
		panic(err)
	}
	if t.Status {
		return "ban"
	}
	return "clean"
}

func SendLkQuery(api string, vals url.Values) string {
	resp, err := http.PostForm(LKAPI+api, vals)
	if err != nil {
		runExtensionCallback(C.CString("secExt"), C.CString("error"), C.CString("SendLkQuery | "+err.Error()))
		SendSentry(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		SendSentry(fmt.Sprintf("Lk response code: %d Api: %v", resp.StatusCode, api))
		return "error"
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		runExtensionCallback(C.CString("secExt"), C.CString("error"), C.CString("SendLkQuery | "+err.Error()))
		SendSentry(err.Error())
	}
	arg := string(body)

	reArg := regexp.MustCompile(`\\`)
	arg = reArg.ReplaceAllString(arg, ``)

	fmt.Println(arg)
	return arg
}

func GetSqfStartCode(script string) string {
	data := url.Values{
		"key":    {"ASDsadasd1231"},
		"script": {script},
	}
	return SendLkQuery("sqf", data)
}

func CheckBan() string {
	data := url.Values{
		"key": {"ASDsadasd1231"},
		"uid": {GetPlayerUid()},
	}

	type Ban struct {
		Ban    string    `json:"ban"`
		Reason string    `json:"ban_reason"`
		Time   time.Time `json:"ban_time"`
	}

	return SendLkQuery("checkban", data)
}

func WritePlayerHardware(args ...string) string {

	vals := url.Values{
		"key":  {"ASDsadasd1231"},
		"data": {fmt.Sprintf("%v", struct2JSON(args))},
	}
	return SendLkQuery("insertban", vals)
}
