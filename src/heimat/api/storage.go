package api

import (
	"io/ioutil"
	"log"
	"os/user"
)

const fileName = ".heimat_token"
const defaultContent = ""

var filePath = homeFolder() + "/" + fileName

func saveToken(token string) {
	writeFile(filePath, []byte(token))
}

func readToken() string {
	return string(readFile(filePath))
}

func removeToken() {
	writeFile(filePath, []byte(""))
}

func writeFile(path string, content []byte) {

	err := ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
}

func readFile(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte(defaultContent)
	}
	return content
}

func homeFolder() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}
