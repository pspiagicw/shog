package app

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

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

		title, summary, body, _, _ := parseMarkdownDocument(readmeBody)
		if title == "" {
			title = humanizeSlug(slug)
		}
		if summary == "" {
			summary = repoURL
		}

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

	type blogEntry struct {
		page    page
		date    time.Time
		hasDate bool
	}

	blogs := make([]blogEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		raw, err := fs.ReadFile(contentFS, filepath.ToSlash(filepath.Join(dir, entry.Name())))
		if err != nil {
			continue
		}

		title, summary, body, meta, dateValue := parseMarkdownDocument(string(raw))
		slug := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		date, hasDate := parseFrontMatterDate(dateValue)
		blogs = append(blogs, blogEntry{
			page: page{
				Title:   fallback(title, humanizeSlug(slug)),
				Summary: fallback(summary, "Blog post"),
				Content: fallback(body, "Empty blog post."),
				Meta:    meta,
			},
			date:    date,
			hasDate: hasDate,
		})
	}

	sort.SliceStable(blogs, func(i, j int) bool {
		if blogs[i].hasDate && blogs[j].hasDate {
			return blogs[i].date.After(blogs[j].date)
		}
		if blogs[i].hasDate != blogs[j].hasDate {
			return blogs[i].hasDate
		}

		return strings.ToLower(blogs[i].page.Title) < strings.ToLower(blogs[j].page.Title)
	})

	pages := make([]page, 0, len(blogs))
	for _, blog := range blogs {
		pages = append(pages, blog.page)
	}

	return pages
}

func parseMarkdownDocument(raw string) (string, string, string, string, string) {
	title, summary, body, dateValue, tags := parseFrontMatter(raw)
	trimmedBody := strings.TrimSpace(body)
	lines := strings.Split(trimmedBody, "\n")

	if title == "" {
		title = firstMarkdownHeading(lines)
	}

	if summary == "" {
		summary = firstMarkdownParagraph(lines)
	}

	return title, summary, trimmedBody, formatMeta(dateValue, tags), dateValue
}

func parseFrontMatter(raw string) (string, string, string, string, []string) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", "", "", "", nil
	}

	if strings.HasPrefix(trimmed, "+++\n") {
		parts := strings.SplitN(trimmed, "\n+++\n", 2)
		if len(parts) == 2 {
			return frontMatterTitle(parts[0]), frontMatterDescription(parts[0]), parts[1], frontMatterDate(parts[0]), frontMatterTags(parts[0])
		}
	}

	if strings.HasPrefix(trimmed, "---\n") {
		parts := strings.SplitN(trimmed, "\n---\n", 2)
		if len(parts) == 2 {
			return frontMatterTitle(parts[0]), frontMatterDescription(parts[0]), parts[1], frontMatterDate(parts[0]), frontMatterTags(parts[0])
		}
	}

	return "", "", raw, "", nil
}

func frontMatterTitle(frontMatter string) string {
	for _, line := range strings.Split(frontMatter, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "title") {
			return unquoteFrontMatterValue(trimmed)
		}
	}

	return ""
}

func frontMatterDate(frontMatter string) string {
	for _, line := range strings.Split(frontMatter, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "date") {
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

func frontMatterTags(frontMatter string) []string {
	lines := strings.Split(frontMatter, "\n")
	var tags []string
	inTagsBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			inTagsBlock = false
			continue
		}

		if inTagsBlock {
			if strings.HasPrefix(trimmed, "- ") {
				tags = append(tags, strings.Trim(strings.TrimSpace(strings.TrimPrefix(trimmed, "- ")), `"'`))
				continue
			}
			if strings.Contains(trimmed, ":") || strings.Contains(trimmed, "=") {
				inTagsBlock = false
			}
		}

		if strings.HasPrefix(trimmed, "tags") {
			value := frontMatterValue(trimmed)
			if value == "" {
				inTagsBlock = true
				continue
			}
			return parseFrontMatterList(value)
		}
	}

	return compactStrings(tags)
}

func parseFrontMatterList(value string) []string {
	trimmed := strings.TrimSpace(value)
	trimmed = strings.TrimPrefix(trimmed, "[")
	trimmed = strings.TrimSuffix(trimmed, "]")
	if trimmed == "" {
		return nil
	}

	parts := strings.Split(trimmed, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.Trim(strings.TrimSpace(part), `"'`)
		if value != "" {
			values = append(values, value)
		}
	}

	return values
}

func formatMeta(dateValue string, tags []string) string {
	parts := make([]string, 0, 2)
	if dateValue != "" {
		parts = append(parts, dateValue)
	}
	if len(tags) > 0 {
		parts = append(parts, strings.Join(tags, ", "))
	}

	return strings.Join(parts, "  •  ")
}

func parseFrontMatterDate(value string) (time.Time, bool) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
	}

	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, true
		}
	}

	return time.Time{}, false
}

func unquoteFrontMatterValue(line string) string {
	value := frontMatterValue(line)
	return strings.Trim(strings.TrimSpace(value), `"'`)
}

func frontMatterValue(line string) string {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) == 1 {
		parts = strings.SplitN(line, ":", 2)
	}
	if len(parts) != 2 {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

func compactStrings(values []string) []string {
	compact := make([]string, 0, len(values))
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			compact = append(compact, value)
		}
	}

	return compact
}
