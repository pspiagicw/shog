package main

import (
	"github.com/pspiagicw/shog/pkg/argparse"
	"github.com/pspiagicw/shog/pkg/ssh"
)

var VERSION string

func main() {
	args := argparse.ParseFlags(VERSION)
	ssh.ServeSSH(args)
}
