package main

import (
	"fmt"
	"os"
)

func debugMsg(msg string) {
	flag, isPresent := os.LookupEnv("LOGLEVEL")
	if isPresent && flag == "1" {
		fmt.Printf(msg)
	}
}

func logMsg(msg string) {
	fmt.Printf(msg)
}
