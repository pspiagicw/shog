package app

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type config struct {
	Title    string       `toml:"title"`
	About    aboutSection `toml:"about"`
	Projects []configPage `toml:"projects"`
	Blog     configPage   `toml:"blog"`
}

type aboutSection struct {
	Title   string `toml:"title"`
	Summary string `toml:"summary"`
}

type configPage struct {
	Title   string `toml:"title"`
	Summary string `toml:"summary"`
	Body    string `toml:"body"`
}

type content struct {
	Title    string
	About    page
	Projects []page
	Blog     page
}

func loadContent(root string) content {
	cfg := defaultConfig()

	if loaded, err := loadConfig(filepath.Join(root, "app.toml")); err == nil {
		cfg = loaded
	}

	aboutBody := defaultAboutBody()
	if loaded, err := os.ReadFile(filepath.Join(root, "about.md")); err == nil {
		aboutBody = strings.TrimSpace(string(loaded))
	}

	projects := make([]page, 0, len(cfg.Projects))
	for _, entry := range cfg.Projects {
		projects = append(projects, page{
			Title:   fallback(entry.Title, "Untitled Project"),
			Summary: fallback(entry.Summary, "No summary provided."),
			Content: fallback(entry.Body, "Add project details in app.toml."),
		})
	}

	if len(projects) == 0 {
		projects = append(projects, page{
			Title:   "Projects",
			Summary: "No projects configured yet.",
			Content: "Add one or more [[projects]] entries in app.toml to populate this tab.",
		})
	}

	return content{
		Title: fallback(cfg.Title, "shog"),
		About: page{
			Title:   fallback(cfg.About.Title, "About"),
			Summary: fallback(cfg.About.Summary, "A short intro delivered over SSH."),
			Content: fallback(aboutBody, defaultAboutBody()),
		},
		Projects: projects,
		Blog: page{
			Title:   fallback(cfg.Blog.Title, "Blog"),
			Summary: fallback(cfg.Blog.Summary, "Writing will show up here later."),
			Content: fallback(cfg.Blog.Body, "No posts yet."),
		},
	}
}

func loadConfig(path string) (config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return config{}, err
	}

	var cfg config
	if err := toml.Unmarshal(raw, &cfg); err != nil {
		return config{}, err
	}

	return cfg, nil
}

func defaultConfig() config {
	return config{
		Title: "shog",
		About: aboutSection{
			Title:   "About",
			Summary: "A short intro delivered over SSH.",
		},
		Projects: []configPage{
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
		Blog: configPage{
			Title:   "Blog",
			Summary: "Writing will show up here later.",
			Body:    "No posts yet.",
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
