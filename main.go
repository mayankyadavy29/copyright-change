package main

import (
	"os"
	"strings"
)

func main() {
	//fmt.Println("Checking/Updating copyright statement !!")
	files := os.Args[1]
	ar := strings.Fields(files)
	for _, f := range ar {
		UpdateCopyright(f)
	}
}
