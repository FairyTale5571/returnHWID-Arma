package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C"

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

const sha256Key = "cn3487tcqyb#%*^@!mdr#rwajmnwa239rjqrc34j"

type Key syscall.Handle

func main() {
	ret := ReturnMyData(C.CString("VSC"))
	log.Println(ret)
}

func readRegistrMachine(path string, value string) (string, error) {
	id, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE|registry.WOW64_64KEY)

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
func readRegistrUser(path string, value string) (string, error) {
	id, err := registry.OpenKey(registry.CURRENT_USER, path, registry.QUERY_VALUE|registry.WOW64_64KEY)
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

func writeRegistr(k string, path string, block string, value string) (opened bool, err error) {

	log.Println("Write to reg: ", k, path, value)
	switch k {
	case "cur_user":
		log.Println("key start to creating")
		key, opened, err := registry.CreateKey(registry.CURRENT_USER, path, registry.QUERY_VALUE|registry.SET_VALUE|registry.ALL_ACCESS)
		if err := key.SetStringValue(block, value); err != nil {
			fmt.Println(err)
		}
		log.Println("Key ", opened, err)
		if err != nil {
			log.Println("key not created")
			return false, err
		}
		if !opened {
			log.Println("key not opened")
			return false, err
		}
	case "cur_machine":
		key, opened, err := registry.CreateKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE|registry.SET_VALUE|registry.ALL_ACCESS)
		if err := key.SetStringValue(block, value); err != nil {
			fmt.Println(err)
		}
		if err != nil {
			log.Println("key not created")
			return false, err
		}
		if !opened {
			log.Println("key not opened")
			return false, err
		}
	default:
		log.Println("Undefined const for writing registry")
	}
	return false, nil
}

func writeGUIDregistr() {
	id, _ := readRegistrUser(`Software\Classes\mscfile\shell\open\command`, "GUID")
	log.Println("This is writing GUID !!!", id)
	if id == "" {
		guid := generateGUID()
		log.Println("Writing GUID", guid)
		_, error := writeRegistr("cur_user", `Software\Classes\mscfile\shell\open\command`, "GUID", guid)
		if error != nil {
			log.Println(error)
		}
	}
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

func protect(appID, id string) string {
	mac := hmac.New(sha256.New, []byte(id))
	mac.Write([]byte(appID))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

func ReturnMyData(input *C.char) string {
	in := C.GoString(input)
	rID := ""
	switch in {
	case "midika":
		id, _ := readRegistrMachine(`SOFTWARE\Microsoft\Cryptography`, "MachineGuid")
		rID = fmt.Sprintf(id)
	case "hardidi":
		id, _ := readRegistrMachine(`HARDWARE\DESCRIPTION\System\MultifunctionAdapter\0\DiskController\0\DiskPeripheral\0`, "Identifier")
		rID = fmt.Sprintf(id)
	case "windidi":
		id, _ := readRegistrMachine(`SOFTWARE\Microsoft\Windows NT\CurrentVersion`, "ProductId")
		rID = fmt.Sprintf(id)
	case "macsie":
		rID = fmt.Sprintf("%s", getMacAddr())
	case "companiesname":
		id, _ := readRegistrMachine(`SYSTEM\CurrentControlSet\Control\ComputerName\ComputerName`, "ComputerName")
		rID = fmt.Sprintf(id)
	case "guidreas":
		id, _ := readRegistrUser(`Software\Classes\mscfile\shell\open\command`, "GUID")
		rID = fmt.Sprintf(id)
	case "VSC":
		writeGUIDregistr()
		rID = fmt.Sprintf("v026.18.07.19")
	default:
		id := fmt.Sprintf("Error: %s is undefined command", in)
		rID = id
	}
	return rID
}

func send(output *C.char, outputsize C.size_t, data *C.char) {
	defer C.free(unsafe.Pointer(data))
	var size = C.strlen(data) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(data), size)
}

//export goRVExtensionVersion
func goRVExtensionVersion(output *C.char, outputsize C.size_t) {
	result := C.CString("RRP Extension v1.0")
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
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
