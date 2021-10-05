package main

import "C"
import (
	"bytes"
	"fmt"
	"gopkg.in/toast.v1"
	"unicode"
)

func CleanStr(str string) string {
	var buf bytes.Buffer
	for _, r := range str {
		if unicode.IsControl(r) {
			fmt.Fprintf(&buf, "\\u%04X", r)
		} else {
			fmt.Fprintf(&buf, "%c", r)
		}
	}
	fmt.Printf(buf.String())
	return buf.String()
}

func SendToast(args ...string) string {

	notification := toast.Notification{
		AppID: CleanStr(args[0]),
		Title: CleanStr(args[1]),
		Message: CleanStr(args[2]),
		Icon: CleanStr(args[3]),
		Loop: false,
		Duration: "short",
	}
	err := notification.Push()
	if err != nil {
		runExtensionCallback(C.CString("returnHWID"), C.CString("error"), C.CString("WinToast | "+err.Error()))
	}
	return "Sended"
}
