package main

import "flag"

type Option struct {
	Module string
	HTTP   bool
	RPC    bool
}

func parse() *Option {
	var opt = new(Option)
	flag.StringVar(&opt.Module, "m", "", "")
	flag.BoolVar(&opt.HTTP, "http", false, "")
	flag.BoolVar(&opt.RPC, "rpc", false, "")
	flag.Parse()
	if !opt.HTTP && !opt.RPC {
		opt.HTTP = true
	}

	return opt

}

func main() {

}
