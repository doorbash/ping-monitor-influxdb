package main

import (
	"log"
	"regexp"
	"runtime"
)

type Constants struct {
	patternNix *regexp.Regexp
	patternWin *regexp.Regexp
	goos       int
}

const (
	WIN = 1
	NIX = 2
)

var constants Constants

func initConstants() {
	var err error
	constants.patternNix, err = regexp.Compile("(.*) bytes from (.*): icmp_seq=(.*) ttl=(.*) time=(.*) ms")
	if err != nil {
		log.Fatalln(err)
	}
	constants.patternWin, err = regexp.Compile("Reply from (.*): bytes=(.*) time[=<](.*)ms TTL=(.*)")
	if err != nil {
		log.Fatalln(err)
	}
	switch _os := runtime.GOOS; _os {
	case "linux":
		constants.goos = NIX
	case "windows":
		constants.goos = WIN
	default:
		log.Fatalln("OS not supported")
	}
}
