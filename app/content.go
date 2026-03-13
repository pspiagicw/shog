package app

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

type config struct {
	Title string       `toml:"title"`
	About aboutSection `toml:"about"`
	Blog  configPage   `toml:"blog"`
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

func loadContent(contentFS fs.FS) content {
	cfg := defaultConfig()

	if loaded, err := loadConfig(contentFS, "app.toml"); err == nil {
		cfg = loaded
	}

	aboutBody := defaultAboutBody()
	if loaded, err := fs.ReadFile(contentFS, "about.md"); err == nil {
		aboutBody = strings.TrimSpace(string(loaded))
	}

	projectURLs := loadProjectURLs(contentFS, "projects.urls")
	projects := make([]page, 0, len(projectURLs))
	for _, repoURL := range projectURLs {
		slug := projectSlug(repoURL)
		readmePath := filepath.ToSlash(filepath.Join("projects", slug+".md"))
		readmeBody := missingREADME(repoURL)
		if loaded, err := fs.ReadFile(contentFS, readmePath); err == nil {
			readmeBody = strings.TrimSpace(string(loaded))
		}

		title, summary, body := projectPageContent(slug, repoURL, readmeBody)
		projects = append(projects, page{
			Title:   title,
			Summary: summary,
			Content: body,
		})
	}

	if len(projects) == 0 {
		projects = append(projects, page{
			Title:   "Projects",
			Summary: "No projects have been fetched yet.",
			Content: missingREADME("Add repository URLs to content/projects.urls and run scripts/fetch-project-readmes.sh before building."),
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

func loadConfig(contentFS fs.FS, filePath string) (config, error) {
	raw, err := fs.ReadFile(contentFS, filePath)
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
		"Use the tab controls to switch between About, Projects, and Blog.",
	}, "\n")
}

func loadProjectURLs(contentFS fs.FS, filePath string) []string {
	raw, err := fs.ReadFile(contentFS, filePath)
	if err != nil {
		return nil
	}

	lines := strings.Split(string(raw), "\n")
	urls := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		urls = append(urls, trimmed)
	}

	return urls
}

func projectSlug(repoURL string) string {
	trimmed := strings.TrimSpace(repoURL)
	trimmed = strings.TrimPrefix(trimmed, "https://")
	trimmed = strings.TrimPrefix(trimmed, "http://")
	trimmed = strings.TrimSuffix(trimmed, "/")

	return filepath.Base(trimmed)
}

func projectPageContent(slug, repoURL, readme string) (string, string, string) {
	trimmed := strings.TrimSpace(readme)
	if trimmed == "" {
		return humanizeSlug(slug), repoURL, missingREADME(repoURL)
	}

	lines := strings.Split(trimmed, "\n")
	title := humanizeSlug(slug)
	if heading := firstMarkdownHeading(lines); heading != "" {
		title = heading
	}

	summary := repoURL
	if excerpt := firstMarkdownParagraph(lines); excerpt != "" {
		summary = excerpt
	}

	return title, summary, trimmed
}

func firstMarkdownHeading(lines []string) string {
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			return strings.TrimSpace(strings.TrimLeft(trimmed, "#"))
		}
	}

	return ""
}

func firstMarkdownParagraph(lines []string) string {
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		return trimmed
	}

	return ""
}

func humanizeSlug(slug string) string {
	replacer := strings.NewReplacer("-", " ", "_", " ")
	value := replacer.Replace(strings.TrimSpace(slug))
	if value == "" {
		return "Project"
	}

	wordPattern := regexp.MustCompile(`\b\w`)
	return wordPattern.ReplaceAllStringFunc(value, strings.ToUpper)
}

func missingREADME(source string) string {
	return strings.Join([]string{
		"README content has not been fetched for this project yet.",
		"",
		"Run scripts/fetch-project-readmes.sh before building to vendor repository READMEs into the binary.",
		"",
		source,
	}, "\n")
}

func fallback(value, defaultValue string) string {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}

	return value
}
