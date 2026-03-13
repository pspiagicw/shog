package main

import (
	"embed"
	"io/fs"
	"os"

	"github.com/pspiagicw/shog/app"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
)

//go:embed content
var contentFiles embed.FS

func main() {
	host := envOrDefault("SHOG_HOST", "0.0.0.0")
	port := envOrDefault("SHOG_PORT", "23234")
	contentFS, err := fs.Sub(contentFiles, "content")
	if err != nil {
		log.Fatal("failed to load embedded content", "err", err)
	}

	server, err := wish.NewServer(
		wish.WithAddress(host+":"+port),
		wish.WithHostKeyPath(".wish_ed25519"),
		wish.WithMiddleware(
			bm.Middleware(app.NewModel(contentFS)),
		),
	)
	if err != nil {
		log.Fatal("failed to create SSH server", "err", err)
	}

	log.Info("starting shog", "address", host+":"+port)
	if err := server.ListenAndServe(); err != nil {
		log.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
