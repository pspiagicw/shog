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
	Blogs    []page
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

	blogs := loadBlogPages(contentFS, "blogs")
	if len(blogs) == 0 {
		blogs = append(blogs, page{
			Title:   fallback(cfg.Blog.Title, "Blog"),
			Summary: fallback(cfg.Blog.Summary, "Writing will show up here later."),
			Content: fallback(cfg.Blog.Body, "No posts yet."),
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
		Blogs:    blogs,
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
	_, _, trimmed := parseMarkdownDocument(readme)
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

func loadBlogPages(contentFS fs.FS, dir string) []page {
	entries, err := fs.ReadDir(contentFS, dir)
	if err != nil {
		return nil
	}

	blogs := make([]page, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		raw, err := fs.ReadFile(contentFS, filepath.ToSlash(filepath.Join(dir, entry.Name())))
		if err != nil {
			continue
		}

		title, summary, body := parseMarkdownDocument(string(raw))
		slug := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		blogs = append(blogs, page{
			Title:   fallback(title, humanizeSlug(slug)),
			Summary: fallback(summary, "Blog post"),
			Content: fallback(body, "Empty blog post."),
		})
	}

	return blogs
}

func parseMarkdownDocument(raw string) (string, string, string) {
	title, summary, body := parseFrontMatter(raw)
	trimmedBody := strings.TrimSpace(body)
	lines := strings.Split(trimmedBody, "\n")

	if title == "" {
		title = firstMarkdownHeading(lines)
	}

	if summary == "" {
		summary = firstMarkdownParagraph(lines)
	}

	return title, summary, trimmedBody
}

func parseFrontMatter(raw string) (string, string, string) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", "", ""
	}

	if strings.HasPrefix(trimmed, "+++\n") {
		parts := strings.SplitN(trimmed, "\n+++\n", 2)
		if len(parts) == 2 {
			return frontMatterFields(parts[0]), frontMatterDescription(parts[0]), parts[1]
		}
	}

	if strings.HasPrefix(trimmed, "---\n") {
		parts := strings.SplitN(trimmed, "\n---\n", 2)
		if len(parts) == 2 {
			return frontMatterFields(parts[0]), frontMatterDescription(parts[0]), parts[1]
		}
	}

	return "", "", raw
}

func frontMatterFields(frontMatter string) string {
	for _, line := range strings.Split(frontMatter, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "title") {
			return unquoteFrontMatterValue(trimmed)
		}
	}

	return ""
}

func frontMatterDescription(frontMatter string) string {
	for _, line := range strings.Split(frontMatter, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "description") || strings.HasPrefix(trimmed, "summary") {
			return unquoteFrontMatterValue(trimmed)
		}
	}

	return ""
}

func unquoteFrontMatterValue(line string) string {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) == 1 {
		parts = strings.SplitN(line, ":", 2)
	}
	if len(parts) != 2 {
		return ""
	}

	return strings.Trim(strings.TrimSpace(parts[1]), `"'`)
}
