package argparse

import (
	"flag"
	"os"

	"github.com/pspiagicw/shog/pkg/help"
)

type Args struct {
	ContentDir string
	Host       string
	Port       string
}

func ParseFlags(version string) *Args {
	args := new(Args)
	flag.StringVar(&args.ContentDir, "content-dir", "content", "Where to load the content from")
	flag.StringVar(&args.Host, "host", "localhost", "Host to bind the SSH server to")
    flag.StringVar(&args.Port, "port", ":2323", "Port to bind the SSH server to")
	showVersion := flag.Bool("version", false, "Show version info")

	flag.Parse()

	if *showVersion {
		help.HelpVersion(version)
		os.Exit(1)
	}

	return args
}
