package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include "extensionCallback.h"
*/
import "C"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"time"
	"unsafe"

	"github.com/rdegges/go-ipify"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func main() {}

var extensionCallbackFnc C.extensionCallback

func runExtensionCallback(name *C.char, function *C.char, data *C.char) C.int {
	return C.runExtensionCallback(extensionCallbackFnc, name, function, data)
}

func callBackArma(input string) {

	id := returnMyData(input, nil)
	temp := (fmt.Sprintf("%s", id))

	name := C.CString("secExt")
	defer C.free(unsafe.Pointer(name))
	function := C.CString(input)
	defer C.free(unsafe.Pointer(function))
	// Make a callback to Arma
	param := C.CString(temp)
	defer C.free(unsafe.Pointer(param))
	runExtensionCallback(name, function, param)
}

//export goRVExtensionRegisterCallback
func goRVExtensionRegisterCallback(fnc C.extensionCallback) {
	extensionCallbackFnc = fnc
}

func returnMyData(input string, errors error) string {

	switch input {
	case "goarch":
		return runtime.GOARCH
	case "checkDRM":
		if !DRMUnlocked {
			return "NO"
		}
		return "YES"
	case "getPlayerUID":
		return GetPlayerUid()
	case "checkInfiBan":
		return CheckInfiBan()
	case "isBan":
		return fmt.Sprintf("%s", CheckBan())
	case "close":
		os.Exit(1)
		return "closed"
	case "loadSqf":
		return GetSqfStartCode()
	case "panic":
		return "panic"
	case "version":
		writeGUIDregistr()
		return fmt.Sprintf("v5 06.10.21")
	case "errors":
		return fmt.Sprintf("Error '%s'", errors)
	case "info":
		return fmt.Sprintf("Created by FairyTale5571. Commands not available for public")
	case "4_c": // clean temp dir
		cleanTemp()
		return fmt.Sprintf("Success")
	case "4_c_a":
		cleanOldFiles()
		return "success"
	}

	if !DRMUnlocked {
		return "U need to disable DRM, contact me FairyTale#5571"
	}

	switch input {
	case "isAdmin":
		return fmt.Sprintf("%d", isAdmin())
	case "get_HWID": // hwid
		id, _ := readReg("local_machine", `SOFTWARE\Microsoft\Cryptography`, "MachineGuid")
		return fmt.Sprintf(id)
	case "get_HDDUID": // HDD_UID
		id, _ := readReg("local_machine", `HARDWARE\DESCRIPTION\System\MultifunctionAdapter\0\DiskController\0\DiskPeripheral\0`, "Identifier")
		return fmt.Sprintf(id)
	case "get_Product": // Product_Win
		return getProductId()
	case "get_Process": // processList
		return fmt.Sprintf(getProcesses())
	case "get_MAC": // MAC
		return fmt.Sprintf(getMacAddr())
	case "get_GUID": // GUID
		id, _ := readReg("current_user", `Software\Classes\mscfile\shell\open\command`, "GUID")
		return fmt.Sprintf(id)
	case "get_IP":
		ip, err := ipify.GetIp()
		if err != nil {
			return "Cant get ip"
		}
		return fmt.Sprintf(ip)
	case "get_GeoIP":
		return getGeoIp()
	case "get_Sd":
		return GetDsName()
	/*************************************/
	case "GetCPU_id":
		return getCpuId()
	case "GetCPU_name":
		return getCpuName()
	/*************************************/
	case "GetMother_id":
		return getMotherId()
	case "GetMother_name":
		return getMotherName()
	/*************************************/
	case "GetBios_id":
		return getBiosId()
	case "GetBios_ReleaseDate":
		return getBiosReleaseDate()
	case "GetBios_Version":
		return getBiosVersion()
	/*************************************/
	case "GetRam_serialNumber":
		return getRamSerialNumber()
	case "GetRam_capacity":
		return getRamCapacity()
	case "GetRam_partNumber":
		return getRamPartNumber()
	case "GetRam_Name":
		return getRamName()
	/*************************************/
	case "GetProduct_Date":
		return getProductInstallDate()
	case "GetProduct_Version":
		return getProductVersion()
	/*************************************/
	case "GetPC_name":
		return getPcName()
	case "Get_SID":
		return getSID()
	case "Get_VRAM_name":
		return getVRAM()
	/*************************************/
	default:
		return fmt.Sprintf("Error '%s' is undefined command, contact Discord FairyTale#5571 for more information", input)
	}

}

func screenCallBack(p_name string, p_uid string) {

	name := C.CString("secExt")
	defer C.free(unsafe.Pointer(name))
	function := C.CString("3_c")
	defer C.free(unsafe.Pointer(function))
	// Make a callback to Arma

	t := time.Now()
	path := fmt.Sprintf("/screenshots/%s_%s/%d/%d/%d", p_name[1:len(p_name)-1], p_uid[1:len(p_uid)-1], t.Year(), t.Month(), t.Day())
	makeScreenshot(path)

	param := C.CString(path)
	defer C.free(unsafe.Pointer(param))

	runExtensionCallback(name, function, param)
	return
}

//export goRVExtensionVersion
func goRVExtensionVersion(output *C.char, outputsize C.size_t) {
	result := C.CString("secExt_v5")
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
}

//export goRVExtension
func goRVExtension(output *C.char, outputsize C.size_t, input *C.char) {

	id := returnMyData(C.GoString(input), nil)
	temp := (fmt.Sprintf("%s", id))
	printInArma(output, outputsize, temp)
}

//export goRVExtensionArgs
func goRVExtensionArgs(output *C.char, outputsize C.size_t, input *C.char, argv **C.char, argc C.int) {
	offset := unsafe.Sizeof(uintptr(0))
	action := C.GoString(input)
	clearArgs := cleanInput(argv, int(argc))

	if f, ok := Cmd2Action[action]; ok {
		printInArma(output, outputsize, f(clearArgs...))
		return
	}

	if !DRMUnlocked {
		printInArma(output, outputsize, "U need to disable DRM, contact me FairyTale#5571")
		return
	}

	switch action {
	case "1_c": // credentials
		var err error
		_creds_json := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv))))
		creds_json := C.GoString(*_creds_json)
		creds_json = creds_json[1 : len(creds_json)-1]

		reToken := regexp.MustCompile(`""`)
		creds_json = reToken.ReplaceAllString(creds_json, `"`)

		b := []byte(creds_json)
		config, err = google.ConfigFromJSON(b, drive.DriveScope)
		fmt.Println(creds_json)
		if err != nil {
			printInArma(output, outputsize, err.Error())
			return
		}
		printInArma(output, outputsize, "Creds accepted")
		return
	case "2_c": // token
		_token_json := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv))))
		token_json := C.GoString(*_token_json)
		token_json = token_json[1 : len(token_json)-1]

		reToken := regexp.MustCompile(`""`)
		token_json = reToken.ReplaceAllString(token_json, `"`)

		fmt.Println(token_json)

		tokenR := bytes.NewReader([]byte(token_json))

		tok = &oauth2.Token{}
		err := json.NewDecoder(tokenR).Decode(tok)
		if err != nil {
			printInArma(output, outputsize, err.Error())
			return
		}
		printInArma(output, outputsize, "Token accepted")
		return
	case "3_c_t":
		_name := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv))))
		_uid := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + offset))

		name := C.GoString(*_name)
		uid := C.GoString(*_uid)

		path := fmt.Sprintf("/%s_%s", name[1:len(name)-1], uid[1:len(uid)-1])
		time.Sleep(5 * time.Second)
		makeScreenshot(path)
		printInArma(output, outputsize, "Done")
		return
	case "3_c": // doit

		_name := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv))))
		_uid := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + offset))

		name := C.GoString(*_name)
		uid := C.GoString(*_uid)

		if extensionCallbackFnc != nil {
			screenCallBack(name, uid)
		} else {

			path := fmt.Sprintf("/%s_%s", name[1:len(name)-1], uid[1:len(uid)-1])
			time.Sleep(5 * time.Second)
			makeScreenshot(path)
			printInArma(output, outputsize, "Done")
		}
		return
	case "1_r": // write_reg
		printInArma(output, outputsize, writeReg(clearArgs[0], clearArgs[1], clearArgs[2], clearArgs[3]))
		return
	case "2_r": // read_reg
		r, _ := readReg(clearArgs[0], clearArgs[1], clearArgs[2])
		printInArma(output, outputsize, r)
		return
	case "3_r": // del_reg
		printInArma(output, outputsize, delReg(clearArgs[0], clearArgs[1], clearArgs[2]))
		return
	case "1_f": // write_file
		printInArma(output, outputsize, writeFile(clearArgs[0], clearArgs[1]))
		return
	case "2_f": // read_file
		printInArma(output, outputsize, readFile(clearArgs[0]))
		return
	case "3_f": // delete_file
		printInArma(output, outputsize, delFile(clearArgs[0]))
		return
	default:
		temp := fmt.Sprintf("Undefined '%s' command", action)
		printInArma(output, outputsize, temp)
		return
	}
}
