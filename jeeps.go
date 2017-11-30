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
	return fmt.Sprintf("pid=%d main=%s args=%s", p.pid, p.mainClass, p.args)
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
	cmd := exec.Command("jps", "-l", "-v", "-m")
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
			// strings.Fields() does not work for finding the main class.
			fields := strings.Split(line, " ")
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
