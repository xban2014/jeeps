package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func main() {
	jps, err := listProcs()
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

func listProcs() ([]JavaProcess, error) {
	cmd := exec.Command("jps", "-l", "-v", "-m")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return parseProcessList(out.String()), nil
}

func parseProcessList(out string) (result []JavaProcess) {
	result = make([]JavaProcess, 0, 5)
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
			result = append(result, NewJavaProcess(pid, mainClass, args))
		}
	}
	return result
}

// JavaProcess - collection of details about a Java process
type JavaProcess struct {
	pid       int
	mainClass string
	args      []string // all args
	xargs     []string // -X<something> args
	sysargs   []string // -D<something> args
	pargs     []string // regular process args
}

// NewJavaProcess constructor
func NewJavaProcess(pid int, mainClass string, args []string) JavaProcess {
	var jp JavaProcess
	jp.pid = pid
	jp.mainClass = mainClass
	jp.args = args
	jp.xargs, jp.sysargs, jp.pargs = splitJavaArgs(args)
	return jp
}

func (p *JavaProcess) plainText() string {
	return fmt.Sprintf(
		"pid=%d main=%s xArgs=%s, sysArgs=%s, prgArgs=%s",
		p.pid, p.mainClass, p.xargs, p.sysargs, p.pargs)
}

func (p *JavaProcess) colorText(dark bool) string {
	return fmt.Sprintf(
		"%s %s %s %s %s",
		colorPid(dark)("%d", p.pid),
		colorMainClass(dark)("%s", p.getMainClass()),
		colorXArgs(dark)("%s", strings.Join(p.xargs, " ")),
		colorSysArgs(dark)("%s", strings.Join(p.sysargs, " ")),
		colorPrgArgs(dark)("%s", strings.Join(p.pargs, " ")))
}

func (p *JavaProcess) getMainClass() string {
	if len(p.mainClass) > 0 {
		return p.mainClass
	}
	return "?"
}

func colorPid(dark bool) func(string, ...interface{}) string {
	return color.New(color.Bold).SprintfFunc()
}

func colorMainClass(dark bool) func(string, ...interface{}) string {
	return color.New(color.FgHiCyan).
		Add(color.Bold).
		Add(color.Underline).
		SprintfFunc()
}

func colorXArgs(dark bool) func(string, ...interface{}) string {
	return color.New(color.BgHiGreen).
		Add(color.FgBlack).
		Add(color.Italic).
		SprintfFunc()
}

func colorSysArgs(dark bool) func(string, ...interface{}) string {
	return color.New(color.BgHiYellow).
		Add(color.FgBlack).
		Add(color.Italic).
		SprintfFunc()
}

func colorPrgArgs(dark bool) func(string, ...interface{}) string {
	return color.New(color.Bold).SprintfFunc()
}

func splitJavaArgs(args []string) (x []string, d []string, p []string) {
	x = make([]string, 0, 5)
	d = make([]string, 0, 5)
	p = make([]string, 0, 5)
	for _, a := range args {
		switch {
		case strings.HasPrefix(a, "-X"):
			x = append(x, a)
		case strings.HasPrefix(a, "-D"):
			d = append(d, a)
		default:
			p = append(p, a)
		}
	}
	return
}
