package main

import (
	"log"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Config string `short:"c" long:"config" description:"path to config file" required:"true"`
}

var opts Options

func main() {
	parser := flags.NewParser(&opts, flags.Default)

	_, err := parser.Parse()

	if err != nil {
		return
	}

	initConstants()

	config := NewConfigParser(opts.Config)
	appConfig, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	InfluxDBHelper := NewInfluxDBHelper(appConfig)

	pingChan := make(chan PingResult, 100)

	for _, host := range appConfig.Hosts {
		NewPingProcess(host, &pingChan).StartPingProcess()
	}

	for ping := range pingChan {
		InfluxDBHelper.SendData(ping)
	}
}
