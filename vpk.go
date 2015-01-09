package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	default:
		path, name := filepath.Split(args[0])
		prevDir, _ := filepath.Abs(".")
		os.Chdir(path)
		defer os.Chdir(prevDir)
		r, err := vpk.OpenReader(name)
		if err != nil {
			panic(err)
		}
		filenames := args[1:]
		var files []*vpk.File
		if len(filenames) == 0 {
			files = r.File
		} else {
			for _, filename := range filenames {
				file := r.FindFile(filename)
				if file != nil {
					files = append(files, file)
				}
			}
		}
		for _, file := range files {
			fp, _ := filepath.Split(file.Name)
			os.MkdirAll(fp, 0755)
			f, err := os.OpenFile(file.Name, os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				fmt.Println(err)
				continue
			}
			defer f.Close()
			rr, err := file.Open()
			if err != nil {
				fmt.Println(err)
				continue
			}
			defer rr.Close()
			_, err = io.Copy(f, rr)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(file.Name)
		}
	}
}
