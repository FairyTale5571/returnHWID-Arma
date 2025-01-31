package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C"

import (
	"fmt"
	"log"
	"time"

	"github.com/StackExchange/wmi"
)

type win32_Processor struct {
	ProcessorId string
	Name        string
	SystemName  string
}

func getCPU() []win32_Processor {

	var dst []win32_Processor
	if err := wmi.Query("SELECT * FROM Win32_Processor", &dst); err != nil {
		SendSentry("getCPU | " + err.Error())
	}
	return dst
}

type win32BaseBoard struct {
	Product      string
	SerialNumber string
}

func getMother() []win32BaseBoard {

	var dst []win32BaseBoard
	if err := wmi.Query("SELECT * FROM Win32_BaseBoard", &dst); err != nil {
		SendSentry("getMother | " + err.Error())
	}
	return dst
}

type Win32_BIOS struct {
	SerialNumber string
	ReleaseDate  time.Time
	Version      string
}

func getBios() []Win32_BIOS {
	var dst []Win32_BIOS
	if err := wmi.Query("SELECT * FROM Win32_BIOS", &dst); err != nil {
		SendSentry("getBios | " + err.Error())
	}
	return dst
}

type Win32PhysicalMemory struct {
	SerialNumber string
	PartNumber   string
	Capacity     uint64
	Manufacturer string
}

func getRAM() []Win32PhysicalMemory {
	var dst []Win32PhysicalMemory
	if err := wmi.Query("SELECT * FROM Win32_PhysicalMemory", &dst); err != nil {
		SendSentry("getRAM | " + err.Error())
	}
	return dst
}

type Win32_OperatingSystem struct {
	Version      string
	InstallDate  time.Time
	SerialNumber string
}

func getOS() []Win32_OperatingSystem {
	var dst []Win32_OperatingSystem
	if err := wmi.Query("SELECT * FROM Win32_OperatingSystem", &dst); err != nil {
		SendSentry("getOS | " + err.Error())
	}
	return dst
}

type Win32_ComputerSystemProduct struct {
	Caption           string
	Description       string
	IdentifyingNumber string
	Name              string
	SKUNumber         string
	Vendor            string
	Version           string
	UUID              string
}

func getCSP() []Win32_ComputerSystemProduct {
	var dst []Win32_ComputerSystemProduct
	if err := wmi.Query("SELECT * FROM Win32_ComputerSystemProduct", &dst); err != nil {
		SendSentry("getCSP | " + err.Error())
	}
	return dst
}

type Win32_VideoController struct {
	Name string
}

func getVRAM() string {
	var dst []Win32_VideoController
	if err := wmi.Query("SELECT * FROM Win32_VideoController", &dst); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("GPU %s", dst[0].Name)
}

type Win32_DiskDrive struct {
	Model        string
	Size         uint64
	SerialNumber string
}

func getDiskDrive() []Win32_DiskDrive {
	var dst []Win32_DiskDrive
	if err := wmi.Query("SELECT * FROM Win32_DiskDrive", &dst); err != nil {
		log.Fatal(err)
	}
	return dst
}
