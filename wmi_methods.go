package main

import (
	"fmt"
	"log"

	"github.com/StackExchange/wmi"
)

func getCpuId() string {
	return fmt.Sprintf(getCPU()[0].ProcessorId)
}

func getCpuName() string {
	return fmt.Sprintf(getCPU()[0].Name)
}

func getMotherId() string {
	return fmt.Sprintf(getMother()[0].SerialNumber)
}

func getMotherName() string {
	return fmt.Sprintf("Motherboard %s", getMother()[0].Product)
}

func getRamSerialNumber() string {
	var numbers []string

	for idx := range getRAM() {
		numbers = append(numbers, getRAM()[idx].SerialNumber)
	}
	return struct2JSON(numbers)
}

func getRamPartNumber() string {
	var numbers []string

	for idx := range getRAM() {
		numbers = append(numbers, getRAM()[idx].PartNumber)
	}
	return struct2JSON(numbers)
}

func getRamName() string {
	var numbers []string

	for idx := range getRAM() {
		numbers = append(numbers, getRAM()[idx].Manufacturer)
	}
	return struct2JSON(numbers)
}

func getRamCapacity() string {
	var memory uint64 = 0
	for idx, _ := range getRAM() {
		memory += getRAM()[idx].Capacity
	}
	size, bytef := ConvertSize(memory)
	return fmt.Sprintf("%d %v", size, bytef)
}

func getProductId() string {
	return fmt.Sprintf(getOS()[0].SerialNumber)
}

func getProductInstallDate() string {
	return fmt.Sprintf("%v", getOS()[0].InstallDate)
}

func getProductVersion() string {
	return fmt.Sprintf(getOS()[0].Version)
}

func getBiosId() string {
	return fmt.Sprintf(getBios()[0].SerialNumber)
}

func getBiosReleaseDate() string {
	return getBios()[0].ReleaseDate.String()
}

func getBiosVersion() string {
	return fmt.Sprintf(getBios()[0].Version)
}

func getPcName() string {
	return fmt.Sprintf(getCPU()[0].SystemName)
}

func getSID() string {
	type Win32_UserAccount struct {
		SID string
	}
	var dst []Win32_UserAccount
	if err := wmi.Query("SELECT * FROM Win32_UserAccount ", &dst); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf(dst[0].SID)
}

func getDiskDrives() string {
	drives := getDiskDrive()

	var drive = "["
	var size uint64
	var str string
	elems := len(drives)
	for idx, _ := range drives {
		size, str = ConvertSize(drives[idx].Size)
		drive += fmt.Sprintf(`["%v %d %v %v"]`, drives[idx].Model, size, str, drives[idx].SerialNumber)
		if elems-1 != idx {
			drive += ","
		}
	}
	drive += "]"
	return drive
}
