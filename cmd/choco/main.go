package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/go-chocolate/chocolate/internal/template"
	"os"
	"path"
	"strings"
)

type Option struct {
	Module string
	Output string
}

func parse() *Option {
	var opt = new(Option)
	flag.StringVar(&opt.Module, "module", "", "")
	flag.StringVar(&opt.Output, "output", ".", "")
	flag.Parse()
	if opt.Output == "" {
		opt.Output = "."
	}

	name := opt.Module
	if n := strings.LastIndex(opt.Module, "/"); n > 0 {
		name = opt.Module[n+1:]
	}
	if name == "" {
		name = "example"
	}
	opt.Output = path.Clean(fmt.Sprintf("%s/%s", opt.Output, name))
	return opt
}

func main() {
	opt := parse()

	files, err := template.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	_, err = os.Stat(opt.Output)
	if err == nil || os.IsExist(err) {
		fmt.Printf("path '%s' has been exists\n", opt.Output)
		os.Exit(1)
	}
	_ = os.MkdirAll(opt.Output, 0777)
	for _, file := range files {
		if err = write(opt, file); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
}

func write(opt *Option, file *template.File) error {
	filename := path.Clean(fmt.Sprintf("%s/%s", opt.Output, file.Filename))
	if file.IsDir {
		return os.Mkdir(filename, 0777)
	}
	var content = file.Content
	if opt.Module != "" {
		content = bytes.ReplaceAll(content, []byte(template.Replacement), []byte(opt.Module))
	}
	return os.WriteFile(filename, content, 0644)
}
