package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	jps, err := runJPS()
	if err != nil {
		log.Fatal(err)
	}
	for _, jp := range jps {
		fmt.Printf("\n%s\n", jp.colorText(true))
	}
}

// there is probably something in the library, I just can't find it.
func lineSeparator() (sep string) {
	if runtime.GOOS == "windows" {
		sep = "\r\n"
	}
	sep = "\n"
	return
}

func listProcs() ([]jpsInfo, error) {
	cmd := exec.Command("jps", "-l", "-v", "-m")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return parseProcessList(out.String()), nil
}

func parseProcessList(out string) (result []jpsInfo) {
	result = make([]jpsInfo, 0, 5)
	lines := strings.Split(out, lineSeparator())
	for _, line := range lines {
		if len(strings.TrimSpace(line)) > 0 {
			// Note: strings.Fields() does not work for finding the main class.
			// It eats away empty fields. IntelliJ Idea shows no main class for instance.
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
