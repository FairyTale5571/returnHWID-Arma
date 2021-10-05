package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"syscall"

	"gopkg.in/yaml.v2"
)

var configPath = getConfigPath()

func LoadServers() []string {
	if err := writeConfig(map[string]interface{}{
		"servers": []string{ // default servers
			"s1.rimasrp.life:8080",
			"dev.rimasrp.life:8080",
			"gtav.rimasrp.life:8080",
		},
	}); err != nil {
		return []string{}
	}
	conf, err := readConfig()
	if err != nil {
		return []string{}
	}

	var servers []string
	for _, server := range conf["servers"].([]interface{}) {
		servers = append(servers, server.(string))
	}
	return servers
}

func WriteServers() error {
	servers := map[string]interface{}{
		"servers": Servers,
	}

	if err := writeConfig(servers); err != nil {
		return err
	}
	return nil
}

func readConfig() (map[string]interface{}, error) {
	f, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Got error while reading %s config file: %v\n", configPath, err))
	}

	var confObj map[string]interface{}
	err = yaml.Unmarshal(f, &confObj)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Got error while unmarshalling %s config file: %v\n", configPath, err))
	}
	return confObj, nil
}

func writeConfig(_data map[string]interface{}) error {
	defer func() {
		// make it hidden
		nameptr, _ := syscall.UTF16PtrFromString(configPath)
		syscall.SetFileAttributes(nameptr, syscall.FILE_ATTRIBUTE_HIDDEN)
	}()

	data, err := yaml.Marshal(_data)
	if err != nil {
		return err
	}

	os.Remove(configPath)
	if err = ioutil.WriteFile(configPath, data, 0777); err != nil {
		return err
	}

	return nil
}

func getConfigPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	sep := string(os.PathSeparator)
	return fmt.Sprint(usr.HomeDir, sep, ".MS.dat")
}
