package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PingProcess struct {
	Address     string
	PingChannel *chan PingResult
}

func (p *PingProcess) StartPingProcess() {
	go func() {
		for {
			var cmd *exec.Cmd
			switch _os := runtime.GOOS; _os {
			case "linux":
				cmd = exec.Command("/usr/bin/ping", p.Address, "-O", "-W", "700")
			case "windows":
				cmd = exec.Command(fmt.Sprintf("%s\\system32\\ping.exe", os.Getenv("WINDIR")), p.Address, "-t", "-w", "700")
			default:
				log.Fatalln("OS not supported")
			}

			stderr, err := cmd.StderrPipe()
			if err != nil {
				log.Println(err)
				continue
			}

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				log.Println(err)
				continue
			}

			if err := cmd.Start(); err != nil {
				log.Println(err)
				continue
			}

			var wg sync.WaitGroup
			wg.Add(2)

			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				scanner := bufio.NewScanner(stderr)
				for scanner.Scan() {
					line := scanner.Text()
					p.parseLine(line)
				}
			}(&wg)

			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					line := scanner.Text()
					p.parseLine(line)
				}
			}(&wg)

			wg.Wait()

			time.Sleep(5 * time.Second)
		}
	}()
}

func (p *PingProcess) parseLine(line string) {
	if strings.Contains(line, "no answer yet") || strings.Contains(line, "Request timed out") {
		*p.PingChannel <- PingResult{
			Time:    -1.0,
			TTL:     -1,
			Address: p.Address,
			Error:   1,
		}
		return
	}

	if strings.Contains(line, "Name or service not known") || strings.Contains(line, "Ping request could not find host") {
		*p.PingChannel <- PingResult{
			Time:    -1.0,
			TTL:     -1,
			Address: p.Address,
			Error:   2,
		}
		return
	}

	patternNix, _ := regexp.Compile("(.*) bytes from (.*): icmp_seq=(.*) ttl=(.*) time=(.*) ms")

	if patternNix.MatchString(line) {
		// NIX
		subMatches := patternNix.FindStringSubmatch(line)
		ttl, err := strconv.Atoi(subMatches[4])
		if err != nil {
			log.Println(err)
			return
		}
		time, err := strconv.ParseFloat(subMatches[5], 32)
		if err != nil {
			log.Println(err)
			return
		}
		*p.PingChannel <- PingResult{
			Time:    time,
			TTL:     ttl,
			Address: p.Address,
			Error:   0,
		}
		return
	}

	patternWin, _ := regexp.Compile("Reply from (.*): bytes=(.*) time[=<](.*)ms TTL=(.*)")

	if patternWin.MatchString(line) {
		// WIN
		subMatches := patternWin.FindStringSubmatch(line)
		ttl, err := strconv.Atoi(subMatches[4])
		if err != nil {
			log.Println(err)
			return
		}
		time, err := strconv.ParseFloat(subMatches[3], 32)
		if err != nil {
			log.Println(err)
			return
		}
		*p.PingChannel <- PingResult{
			Time:    time,
			TTL:     ttl,
			Address: p.Address,
			Error:   0,
		}
	}
}

func NewPingProcess(address string, pingChannel *chan PingResult) *PingProcess {
	return &PingProcess{
		Address:     address,
		PingChannel: pingChannel,
	}
}
