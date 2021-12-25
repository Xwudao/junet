package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func CheckErrWithStatus(err error) {
	if err != nil {
		Error(err)
		os.Exit(0)
	}
}
func LoadFiles(dir string, filter func(filename string) bool) (filenames []string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		filename := filepath.Join(dir, file.Name())
		if file.IsDir() {
			filenames = append(filenames, LoadFiles(filename, filter)...)
		} else {
			if filter(filename) {
				filenames = append(filenames, filename)
			}
		}
	}
	return
}

func Error(err error) {
	if err != nil {
		log.SetPrefix("ERROR")
		log.Println(err.Error())
	}
}
func Info(s string) {
	log.SetPrefix("")
	log.Println(s)
}
