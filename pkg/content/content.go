package content

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/adrg/frontmatter"
	"github.com/pspiagicw/shog/pkg/argparse"
)

func GetSplash(args *argparse.Args) string {
	contents, err := os.ReadFile(args.Splash)
	if err != nil {
		return DEFAULT_SPLASH
	}
	return string(contents)
}

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
	_, err := frontmatter.MustParse(strings.NewReader(content), &matter, []*frontmatter.Format{
		frontmatter.NewFormat("+++", "+++", toml.Unmarshal),
	}...)
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
			return err
		}
		content, err := readFile(file)
		if err != nil {
			return err
		}
		frontMatter := parseFile(content)
		if !info.IsDir() {
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

func readFile(file string) (string, error) {
	contents, err := os.ReadFile(file)
	if err != nil {
		return "", nil
	}
	return string(contents), nil
}
