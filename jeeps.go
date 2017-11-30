package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

const (
	// lineSeparator should be somewhere in a library
	lineSeparator = "\n"
)

// JavaProcess - collection of details about a Java process
type JavaProcess struct {
	pid       int
	mainClass string
	args      []string
}

func (p *JavaProcess) String() string {
	return p.mainClass
}

func main() {
	jps, err := listProcs()
	if err != nil {
		log.Fatal(err)
	}
	for _, jp := range jps {
		fmt.Println(jp.String())
	}
}

func listProcs() ([]JavaProcess, error) {
	// run the jps command - we get a nice error if it cannot be found in $PATH
	cmd := exec.Command("jps", "-l", "-v")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var jps = make([]JavaProcess, 0)

	lines := strings.Split(out.String(), lineSeparator)
	for _, line := range lines {
		if len(strings.TrimSpace(line)) > 0 {
			fields := strings.Fields(line)
			pid, err := strconv.Atoi(fields[0])
			if err != nil {
				log.Fatal("could not convert pid from: " + fields[0])
			}
			mainClass := fields[1]
			args := fields[2:]
			jps = append(jps, JavaProcess{pid, mainClass, args})
		}
	}

	return jps, err
}
