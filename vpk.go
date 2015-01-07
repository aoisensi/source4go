package main

import (
	"flag"
	"fmt"
	"source4go/vpk"
)

func main() {
	flag.Parse()
	args := flag.Args()
	mode := ""
	if len(args) != 0 {
		mode = args[0]
	}
	switch mode {
	case "l":
		path := args[1]
		r, err := vpk.OpenReader(path)
		if err != nil {
			panic(err)
		}
		defer r.Close()
		for _, file := range r.File {
			fmt.Println(file.Name)
		}
	}
}
