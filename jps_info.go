package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type jpsInfo struct {
	pid       int
	mainClass string
	xArgs     []string // -X<something> args
	vmArgs    []string // -D<something> args
	prgArgs   []string // regular process args
}

type argSpec struct {
	key      string
	isPrefix bool
}

var (
	xArgSpecs = []argSpec{
		{"-X", true},
		{"-agentlib:", true},
		{"-ea", false},
	}
	vmArgSpecs = []argSpec{
		{"-D", true},
	}
)

func newJpsInfo(pid int, mainClass string, args []string) jpsInfo {
	var jp jpsInfo
	jp.pid = pid
	jp.mainClass = mainClass
	jp.xArgs, jp.vmArgs, jp.prgArgs = classifyArgs(args)
	return jp
}

func classifyArgs(args []string) (x []string, s []string, p []string) {
	x = make([]string, 0, 5)
	s = make([]string, 0, 5)
	p = make([]string, 0, 5)
	for _, a := range args {
		switch {
		case classifyArg(xArgSpecs, a):
			x = append(x, a)
		case classifyArg(vmArgSpecs, a):
			s = append(s, a)
		default:
			p = append(p, a)
		}
	}
	return
}

func classifyArg(specs []argSpec, arg string) bool {
	for _, spec := range specs {
		if spec.isPrefix && strings.HasPrefix(arg, spec.key) {
			return true
		}
		if !spec.isPrefix && arg == spec.key {
			return true
		}
	}
	return false
}

func (p *jpsInfo) plainText() string {
	return fmt.Sprintf(
		"pid=%d main=%s xArgs=%s, sysArgs=%s, prgArgs=%s",
		p.pid, p.mainClass, p.xArgs, p.vmArgs, p.prgArgs)
}

func (p *jpsInfo) colorText(dark bool) string {
	return fmt.Sprintf(
		"%s %s %s %s %s",
		colorPid(dark)("%d", p.pid),
		colorMainClass(dark)("%s", p.getMainClass()),
		colorXArgs(dark)("%s", strings.Join(p.xArgs, " ")),
		colorSysArgs(dark)("%s", strings.Join(p.vmArgs, " ")),
		colorPrgArgs(dark)("%s", strings.Join(p.prgArgs, " ")))
}

func (p *jpsInfo) getMainClass() string {
	if len(p.mainClass) > 0 {
		return p.mainClass
	}
	return "<main>"
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
