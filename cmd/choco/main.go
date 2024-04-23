package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/go-chocolate/chocolate/internal/args"
	"github.com/go-chocolate/chocolate/internal/template"
	"github.com/go-chocolate/chocolate/internal/version"
)

type option struct {
	template.Option
	command string
}

func parse() *option {
	cmd := args.Parse()
	opt := &option{}
	opt.command = cmd.Command
	opt.Module = cmd.Options.GetAlias([]string{"-module", "module", "m"}, "example")
	opt.Output = cmd.Options.GetAlias([]string{"-output", "output", "o"}, ".")
	name := opt.Module
	if name == "" {
		name = "example"
	}
	if opt.Output == "" {
		opt.Output = "."
	}
	if n := strings.LastIndex(name, "/"); n > 0 {
		name = name[n+1:]
	}
	opt.Output = fmt.Sprintf("%s/%s", opt.Output, name)
	return opt
}

func usage() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("Usage for gokit %s:\n", version.Version))
	buf.WriteString(fmt.Sprintf("  gokit [command] [options]\n"))
	buf.WriteString(fmt.Sprintf("Commands:\n"))
	buf.WriteString(fmt.Sprintf("  create    Create a new project\n"))
	buf.WriteString(fmt.Sprintf("Options:\n"))
	buf.WriteString(fmt.Sprintf("  --module, -m    Module name, default is 'example'.\n"))
	buf.WriteString(fmt.Sprintf("  --output, -o    Output path, default is current path.\n"))
	buf.WriteString("Example:\n  gokit create --module github.com/example/example --output . ")
	return buf.String()
}

func main() {
	opt := parse()
	switch opt.command {
	case "create":
		build(opt)
		fmt.Printf("project '%s' created successfully\n", opt.Module)
		fmt.Printf("run 'cd %s && go run main.go' to quick start.", opt.Output)
	case "", "help":
		fmt.Println(usage())
	}
}

func build(opt *option) {
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
		if err = template.Write(file, opt.Option); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
}
