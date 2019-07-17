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
	"fmt"
	"log"
	"net"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

const sha256Key = "cn3487tcqyb#%*^@!mdr#rwajmnwa239rjqrc34j"

func main() {
	id, _ := readRegistr(`SOFTWARE\Poison`, "GUID")
	fmt.Println(id)
}

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

func writeRegistr(input1 string, value string) (string, error) {
	id, opened, err := registry.CreateKey(registry.LOCAL_MACHINE, input1, registry.QUERY_VALUE|registry.WOW64_64KEY)
	if err != nil {
		return "", err
	}
	defer id.Close()

	s, _, err := id.GetStringsValue(value)
	if err != nil {
		return "", err
	}
	return s, nil
}

func generateGUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

//export goRVExtensionVersion
func goRVExtensionVersion(output *C.char, outputsize C.size_t) {
	result := C.CString("HWID Return v1.0")
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
}

func getMacAddr() (addr string) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				// Don't use random as we have a real address
				addr = i.HardwareAddr.String()
				break
			}
		}
	}
	return
}

func ReturnMyData(input *C.char) string {
	in := C.GoString(input)
	rID := ""
	switch in {
	case "midika":
		id, _ := readRegistr(`SOFTWARE\Microsoft\Cryptography`, "MachineGuid")
		rID = fmt.Sprintf(id)
	case "hardidi":
		id, _ := readRegistr(`HARDWARE\DESCRIPTION\System\MultifunctionAdapter\0\DiskController\0\DiskPeripheral\0`, "Identifier")
		rID = fmt.Sprintf(id)
	case "windidi":
		id, _ := readRegistr(`SOFTWARE\Microsoft\Windows NT\CurrentVersion`, "ProductId")
		rID = fmt.Sprintf(id)
	case "macsie":
		rID = fmt.Sprintf("%s", getMacAddr())
	case "companiesname":
		id, _ := readRegistr(`SYSTEM\CurrentControlSet\Control\ComputerName\ComputerName`, "ComputerName")
		rID = fmt.Sprintf(id)
	case "guidreas":
		rID = generateGUID()
	case "VSC":
		rID = fmt.Sprintf("v025.14.07.19")
	default:
		id := fmt.Sprintf("Error: %s is undefined command", in)
		rID = id
	}
	return rID
}

//export goRVExtension
func goRVExtension(output *C.char, outputsize C.size_t, input *C.char) {
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
