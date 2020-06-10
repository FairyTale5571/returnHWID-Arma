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
	"log"
	"net"
	"unsafe"

	ps "github.com/mitchellh/go-ps"
	"golang.org/x/sys/windows/registry"
)

func main() {}

func getProcesses() string {
	procs, err := ps.Processes()
	if err != nil {
		ReturnMyData("errors", err)
	}

	result := make(map[string]struct{})
	for _, proc := range procs {
		name := proc.Executable()
		if _, ok := result[name]; !ok {
			result[name] = struct{}{}
		}
	}
	names := []string{}
	for key := range result {
		names = append(names, key)
	}
	return fmt.Sprintf("%v\n", struct2JSON(names))
}

func struct2JSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func readRegistrMachine(path string, value string) (string, error) {
	id, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE|registry.WOW64_64KEY)

	if err != nil {
		ReturnMyData("errors", err)
		return "", err
	}
	defer id.Close()

	s, _, err := id.GetStringValue(value)
	if err != nil {
		ReturnMyData("errors", err)
		return "", err
	}
	return s, nil
}

func readRegistrUser(path string, value string) (string, error) {
	id, err := registry.OpenKey(registry.CURRENT_USER, path, registry.QUERY_VALUE|registry.WOW64_64KEY)
	if err != nil {
		ReturnMyData("errors", err)
		return "", err
	}
	defer id.Close()

	s, _, err := id.GetStringValue(value)
	if err != nil {
		ReturnMyData("errors", err)
		return "", err
	}
	return s, nil
}

func writeRegistr(path string, block string, value string) (opened bool, err error) {
	key, opened, err := registry.CreateKey(registry.CURRENT_USER, path, registry.QUERY_VALUE|registry.SET_VALUE|registry.ALL_ACCESS)
	if err := key.SetStringValue(block, value); err != nil {
		fmt.Println(err)
	}
	log.Println("Key ", opened, err)
	if err != nil {
		ReturnMyData("errors", err)
		return false, err
	}
	if !opened {
		ReturnMyData("errors", err)
		return false, err
	}
	return false, nil
}

func writeGUIDregistr() {
	id, _ := readRegistrUser(`Software\Classes\mscfile\shell\open\command`, "GUID")
	if id == "" {
		guid := generateGUID()
		_, error := writeRegistr(`Software\Classes\mscfile\shell\open\command`, "GUID", guid)
		if error != nil {
			ReturnMyData("errors", error)
		}
	}
}

func generateGUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		ReturnMyData("errors", err)
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

//export ReturnMyData
func ReturnMyData(input string, errors error) string {
	rID := ""
	switch input {
	case "processList":
		rID = fmt.Sprintf(getProcesses())
	case "MAC":
		rID = fmt.Sprintf(getMacAddr())
	case "GUID":
		id, _ := readRegistrUser(`Software\Classes\mscfile\shell\open\command`, "GUID")
		rID = fmt.Sprintf(id)
	case "version":
		writeGUIDregistr()
		rID = fmt.Sprintf("v0.28|08.06.20")
	case "errors":
		rID = fmt.Sprintf("Error  '%s'", errors)
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
	result := C.CString("RRPHW v.0.28")
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
}

//export goRVExtension
func goRVExtension(output *C.char, outputsize C.size_t, input *C.char) {
	id := ReturnMyData(C.GoString(input), nil)
	temp := (fmt.Sprintf("%s", id))
	// Return a result to Arma
	result := C.CString(temp)
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
}

//export goRVExtensionArgs
func goRVExtensionArgs(output *C.char, outputsize C.size_t, input *C.char, argv **C.char, argc C.int) {
	var offset = unsafe.Sizeof(uintptr(0))
	var out []string
	for index := C.int(0); index < argc; index++ {
		out = append(out, C.GoString(*argv))
		argv = (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + offset))
	}
	temp := fmt.Sprintf("Function: %s nb params: %d params: %s!", C.GoString(input), argc, out)

	// Return a result to Arma
	result := C.CString(temp)
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
}
