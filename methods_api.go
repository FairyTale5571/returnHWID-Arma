package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include "extensionCallback.h"
*/
import "C"

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
	}
	if t.Status {
		return "ban"
	}
	return "clean"
}

func SendLkQuery(api string, vals url.Values) string {
	resp, err := http.PostForm(LKAPI+api, vals)
	if err != nil {
		SendSentry("SendLkQuery | " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		SendSentry(fmt.Sprintf("Lk response code: %d Api: %v", resp.StatusCode, api))
		return "error"
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		SendSentry("SendLkQuery | " + err.Error())
	}
	arg := string(body)

	return arg
}

func GetSqfStartCode(script string) string {
	data := url.Values{
		"key":    {"ASDsadasd1231"},
		"script": {script},
	}
	ret := SendLkQuery("sqf", data)

	return ret
}

func CheckBan() string {
	data := url.Values{
		"key": {"ASDsadasd1231"},
		"uid": {GetPlayerUid()},
	}

	ret := SendLkQuery("checkban", data)

	t := Ban{}
	err := json.Unmarshal([]byte(ret), &t)
	if err != nil {
		SendSentry(err.Error())
	}

	return fmt.Sprintf(`[%v, "%v"]`, t.Ban, t.Reason)
}

func WritePlayerHardware(args []string) string {

	vals := url.Values{
		"key":  {"ASDsadasd1231"},
		"uid":  {GetPlayerUid()},
		"data": {fmt.Sprintf("%v", struct2JSON(args))},
	}
	return SendLkQuery("inserthardware", vals)
}
