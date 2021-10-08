package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C"

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"syscall"
	"time"
	"unicode"
	"unsafe"

	"github.com/getsentry/sentry-go"
	"golang.org/x/sys/windows/registry"

	"google.golang.org/api/drive/v3"

	"fmt"
	"strings"
)

func printInArma(output *C.char, outputsize C.size_t, input string) {
	result := C.CString(input)
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
}

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
		runExtensionCallback(C.CString("secExt"), C.CString("error"), C.CString("getGoCategory | Unsupported category "))
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

func readReg(category, path, value string) (string, error) {
	goCategory := getGoCategory(category)
	k, err := registry.OpenKey(goCategory, path, registry.QUERY_VALUE|registry.WOW64_64KEY)
	if err != nil {
		return err.Error(), err
	}

	s, _, err := k.GetStringValue(value)
	if err != nil {
		return err.Error(), err
	}
	return s, nil
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
		return fmt.Sprintf("MkDir Error: %s", err.Error())
	}

	if _, err := os.Stat(path); err == nil {
		return "Already exist"
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Sprintf("Create Error: %s", err.Error())
	}
	defer f.Close()

	_, err = f.Write([]byte(data))
	if err != nil {
		return fmt.Sprintf("Write Error: %s", err.Error())
	}

	nameptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		os.Remove(path)
		return fmt.Sprintf("Nameptr Error: %s", err.Error())
	}

	err = syscall.SetFileAttributes(nameptr, syscall.FILE_ATTRIBUTE_HIDDEN)
	if err != nil {
		os.Remove(path)
		return fmt.Sprintf("Attribute Error: %s", err.Error())
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

func InitSentry() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://75311a6f34fd40bbb8cf762330b75eb5@o482351.ingest.sentry.io/5982259",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)
	SentryInit = true
}

func SendSentry(input string) {
	if !SentryInit {
		InitSentry()
	}
	input = fmt.Sprintf("UID: %v | Error: %v", GetPlayerUid(), input)
	sentry.CaptureMessage(input)
}

func SendSetryArma(args ...string) string {
	SendSentry(fmt.Sprintf("%v", struct2JSON(args)))

	return "Sentry"
}
