package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	Title     string       `json:"title"`
	AboutPage aboutPage    `json:"about_page"`
	Pages     []configPage `json:"pages"`
}

type aboutPage struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type configPage struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Body    string `json:"body"`
}

func loadContent(root string) (string, []page) {
	cfg := defaultConfig()

	if loaded, err := loadConfig(filepath.Join(root, "app.config")); err == nil {
		cfg = loaded
	}

	aboutBody := defaultAboutBody()
	if loaded, err := os.ReadFile(filepath.Join(root, "about.md")); err == nil {
		aboutBody = strings.TrimSpace(string(loaded))
	}

	pages := []page{
		{
			Title:   fallback(cfg.AboutPage.Title, "About"),
			Summary: fallback(cfg.AboutPage.Summary, "A short intro delivered over SSH."),
			Content: fallback(aboutBody, defaultAboutBody()),
		},
	}

	for _, entry := range cfg.Pages {
		pages = append(pages, page{
			Title:   entry.Title,
			Summary: entry.Summary,
			Content: entry.Body,
		})
	}

	return fallback(cfg.Title, "shog"), pages
}

func loadConfig(path string) (config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return config{}, err
	}

	var cfg config
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return config{}, err
	}

	return cfg, nil
}

func defaultConfig() config {
	return config{
		Title: "shog",
		AboutPage: aboutPage{
			Title:   "About",
			Summary: "A short intro delivered over SSH.",
		},
		Pages: []configPage{
			{
				Title:   "Project One",
				Summary: "Replace this with a real flagship project.",
				Body: strings.Join([]string{
					"Problem: Describe the problem space in one sentence.",
					"",
					"Build: Note the stack, architecture, or unusual constraint.",
					"",
					"Result: Summarize the outcome, adoption, or what you learned.",
				}, "\n"),
			},
			{
				Title:   "Project Two",
				Summary: "A second page keeps the navigation model honest.",
				Body: strings.Join([]string{
					"Highlight a contrasting project here.",
					"",
					"This can later become a richer card, timeline, or interactive selector.",
					"",
					"A future blog page can slot into the same page list without changing the app shape.",
				}, "\n"),
			},
		},
	}
}

func defaultAboutBody() string {
	return strings.Join([]string{
		"Hi, I'm Your Name.",
		"",
		"I build software, explore ideas, and use this terminal app as a compact portfolio.",
		"",
		"Use the arrow keys or h/l to move between pages.",
	}, "\n")
}

func fallback(value, defaultValue string) string {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}

	return value
}
