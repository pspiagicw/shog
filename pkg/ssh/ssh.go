package ssh

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/shog/pkg/argparse"
	"github.com/pspiagicw/shog/pkg/content"
	"github.com/pspiagicw/shog/pkg/tui"
)

func ServeSSH(args *argparse.Args) {

	// Handle SIGTERM, SIGINT signals graciously
	done := createSignalHandler()
	s := createSerer(args)
	startServer(s, done)
	shutdownServer(s, done)
}
func createSignalHandler() chan os.Signal {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return done
}

func listenAndServe(server *ssh.Server, done chan os.Signal) {
	err := server.ListenAndServe()
	if err != nil && errors.Is(err, ssh.ErrServerClosed) {
		goreland.LogError("Couldn't start the sever: %q", err)
		done <- nil

	}
}
func startServer(server *ssh.Server, done chan os.Signal) {
	goreland.LogInfo("Starting SSH server")
	go listenAndServe(server, done)
}
func createSerer(args *argparse.Args) *ssh.Server {
	blogs := content.GetBlogs(args)
	s, err := wish.NewServer(
		wish.WithAddress("localhost:2323"),
		wish.WithHostKeyPath(".ssh/term_info"),
		wish.WithMiddleware(
			bm.Middleware(tui.EntryGenerator(args, blogs)),
			lm.Middleware(),
		),
	)
	if err != nil {
		goreland.LogFatal("Couldn't start the server: %q", err)
	}
	return s
}
func shutdownServer(server *ssh.Server, done chan os.Signal) {
	<-done

	goreland.LogInfo("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()

	err := server.Shutdown(ctx)
	if err != nil && errors.Is(err, ssh.ErrServerClosed) {
		goreland.LogError("Couldn't stop the server: %v", err)
	}
}
