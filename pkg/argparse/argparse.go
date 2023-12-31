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
	Splash     string
	Name       string
	Slogan     string
}

func ParseFlags(version string) *Args {
	args := new(Args)
	flag.StringVar(&args.ContentDir, "content-dir", "content", "Where to load the content from")
	flag.StringVar(&args.Host, "host", "localhost", "Host to bind the SSH server to")
	flag.StringVar(&args.Port, "port", ":2323", "Port to bind the SSH server to")
	flag.StringVar(&args.Splash, "splash", "splash", "Content to show in the home.")
	flag.StringVar(&args.Name, "name", "shog", "Name of the blog")
	flag.StringVar(&args.Slogan, "slogal", "Serve your blogs over SSH!", "Slogan to show")
	showVersion := flag.Bool("version", false, "Show version info")

	flag.Parse()

	if *showVersion {
		help.HelpVersion(version)
		os.Exit(1)
	}

	return args
}
