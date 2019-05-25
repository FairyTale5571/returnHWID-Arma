package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C"

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

const sha256Key = "cn3487tcqyb#%*^@!mdr#rwajmnwa239rjqrc34j"

func readRegistr(input1 string, value string) (string, error) {
	id, err := registry.OpenKey(registry.LOCAL_MACHINE, input1, registry.QUERY_VALUE|registry.WOW64_64KEY)
	if err != nil {
		return "", err
	}
	defer id.Close()

	s, _, err := id.GetStringValue(value)
	if err != nil {
		return "", err
	}
	return s, nil
}

//export RVExtensionVersion
func RVExtensionVersion(output *C.char, outputsize C.size_t) {
	result := C.CString("HWID Return v1.0")
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
}

func protect(appID, id string) string {
	mac := hmac.New(sha256.New, []byte(id))
	mac.Write([]byte(appID))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

func ReturnMyData(input *C.char) string {
	in := C.GoString(input)
	rID := ""
	switch in {
	case "Machine_ID":
		id, _ := readRegistr(`SOFTWARE\Microsoft\Cryptography`, "MachineGuid")
		rID = fmt.Sprintf(protect(sha256Key, id))
	case "HDD_UID":
		id, _ := readRegistr(`HARDWARE\DESCRIPTION\System\MultifunctionAdapter\0\DiskController\0\DiskPeripheral\0`, "Identifier")
		rID = fmt.Sprintf(protect(sha256Key, id))
	case "Product_Win":
		id, _ := readRegistr(`SOFTWARE\Microsoft\Windows NT\CurrentVersion`, "ProductId")
		rID = fmt.Sprintf(protect(sha256Key, id))
	default:
		id := fmt.Sprintf("Error: %s is undefined command", in)
		rID = id
	}
	return rID
}

//export RVExtension
func RVExtension(output *C.char, outputsize C.size_t, input *C.char) {
	id := ReturnMyData(input)
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

func main() {}
