package main

import (
	"log"

	"github.com/Xwudao/junet/cmd/junet/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
