package main

import (
	"bufio"
	"encoding/json"
	"os/user"
	"os"
	"path/filepath"
)

func jsonify(config *Configuration) (string) {
	byteArray, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	return string(byteArray)
}

func readJson(jsonString string) *Configuration {
	var result *Configuration
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		panic(err)
	}
	return result
}

func handleException(err error) {
	if err != nil {
		panic(err)
	}
}

type Configuration struct {
	FullName string `json:"FullName"`
	GithubUsername string `json:"GithubUsername"`
	DefaultLang string `json:"DefaultLang"`
	DefaultLicense string `json:"DefaultLicense"`
}

func config(fullName, githubUsername, defaultLang, defaultLicense string) {
	// * defining path of hydra config file
	usr, _ := user.Current()
	configFile := filepath.Join(usr.HomeDir, "hydra_config.json")
	_, e := os.Stat(configFile)
	
	// * creating a file in case it doesnt exists
	if e != nil {
		f, err := os.Create(configFile)
		handleException(err)
		defaultConfig := Configuration{FullName: "", GithubUsername: "", DefaultLang: "", DefaultLicense: "MIT"}
		_, er := f.WriteString(jsonify(&defaultConfig))
		handleException(er)
		f.Close()
	}

	// * reading data from the file
	file, ferr := os.Open(configFile)
	handleException(ferr)
	wholeText := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		wholeText = wholeText + line
	}
	file.Close()
	
	// * writing new config to the file by first deleting it
	configStruct := readJson(wholeText)
	configStruct.FullName = fullName
	configStruct.GithubUsername = githubUsername
	configStruct.DefaultLang = defaultLang
	configStruct.DefaultLicense = defaultLicense

	os.Remove(configFile)
	f, err := os.Create(configFile)
	handleException(err)
	_, er := f.WriteString(jsonify(configStruct))
	handleException(er)
	f.Close()

}