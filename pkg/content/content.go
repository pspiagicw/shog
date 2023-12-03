package content

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/pspiagicw/shog/pkg/argparse"
)

type Blog struct {
	Filepath string
	FrontMatter
	Content string
}
type FrontMatter struct {
	BlogTitle  string   `toml:"title"`
	Author     []string `toml:"authors"`
	Tags       []string `toml:"tags"`
	Date       string   `toml:"date"`
	Categories []string `toml:"categories"`
	Draft      bool     `toml:"draft"`
}

func (b Blog) Title() string       { return b.BlogTitle }
func (b Blog) Description() string { return b.Filepath }
func (b Blog) FilterValue() string { return b.BlogTitle }

func parseFile(content string) *FrontMatter {
	var matter FrontMatter
	_, err := frontmatter.Parse(strings.NewReader(content), &matter)
	if err != nil {
		return &matter
	}
	return &matter
}
func GetBlogs(args *argparse.Args) []Blog {
	items := []Blog{}
	fmt.Println(args.ContentDir)
	err := filepath.Walk(args.ContentDir, func(file string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		content := readFile(file)
		frontMatter := parseFile(content)
		fmt.Println(content)
		if !info.IsDir() && frontMatter.BlogTitle != "" {
			items = append(items, Blog{
				FrontMatter: *frontMatter,
				Content:     content,
				Filepath:    file,
			})
		}
		return nil
	})
	if err != nil {
		return []Blog{}
	}
	return items
}

func readFile(file string) string {
	contents, err := os.ReadFile(file)
	if err != nil {
		return ""
	}
	return string(contents)
}
