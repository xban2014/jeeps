package main

import (
	"bytes"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

// can't get ps auxww on MacOS. There is some nasty command that approximates it
// (ps -Ao user,pid,%cpu,%mem,vsz,rss,tt,stat,start,time,command).
func runPSAUXWW() ([]jpsInfo, error) {
	cmd := exec.Command("ps", "auxww")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return parsePSAUXWWOutput(out.String()), nil
}

func parsePSAUXWWOutput(out string) (result []jpsInfo) {
	result = make([]jpsInfo, 0, 5)
	lines := strings.Split(out, lineSeparator())
	for _, line := range lines {
		if len(strings.TrimSpace(line)) > 0 {
			// Note: strings.Fields() does not work for finding the main class.
			// It eats away empty fields. For ex, IntelliJ Idea shows no main class for instance.
			fields := strings.Split(line, " ")
			pid, err := strconv.Atoi(fields[0])
			if err != nil {
				log.Fatal("could not convert pid from: " + fields[0])
			}
			mainClass := fields[1]
			if mainClass == "sun.tools.jps.Jps" {
				continue
			}
			args := fields[2:]
			result = append(result, newJpsInfo(pid, mainClass, args))
		}
	}
	return result
}
