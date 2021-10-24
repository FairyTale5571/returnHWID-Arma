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
	return fmt.Sprintf(getRAM()[0].SerialNumber)
}

func getRamPartNumber() string {
	return fmt.Sprintf(getRAM()[0].PartNumber)
}

func getRamName() string {
	return fmt.Sprintf("RAM %s", getRAM()[0].Manufacturer)
}

func getRamCapacity() string {
	return fmt.Sprintf("%d", getRAM()[0].Capacity)
}

func getProductId() string {
	return fmt.Sprintf(getOS()[0].SerialNumber)
}

func getProductInstallDate() string {
	return fmt.Sprintf(getOS()[0].InstallDate.String())
}
func getProductVersion() string {
	return fmt.Sprintf(getOS()[0].Version)
}

func getBiosId() string {
	return fmt.Sprintf(getBios()[0].SerialNumber)
}

func getBiosReleaseDate() string {
	return fmt.Sprintf(getBios()[0].ReleaseDate.String())
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

	return struct2JSON(drives)
}
