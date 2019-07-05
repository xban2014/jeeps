package main

import (
	"fmt"
	"log"
	"runtime"
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
	} else {
		sep = "\n"
	}
	return
}
