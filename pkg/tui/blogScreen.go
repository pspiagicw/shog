package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/shog/pkg/content"
)

type BlogViewer struct {
	blogViewport viewport.Model
}

func (m Model) viewBlogScreen() string {
	view := strings.Builder{}
	view.WriteString(m.blogViewer.blogViewport.View())
	// view.WriteString(fmt.Sprintf("\n%q\n", ))
	return view.String()
}
func (b *BlogViewer) SetBlog(width int, blog content.Blog) {
	b.blogViewport.SetContent(renderContent(width, blog.Content))
}
func (b *BlogViewer) Update(msg tea.Msg) tea.Cmd {
	v, cmd := b.blogViewport.Update(msg)
	b.blogViewport = v
	return cmd

}
func (b *BlogViewer) SetSize(width, height int) tea.Cmd {
	b.blogViewport.Width = width
	b.blogViewport.Height = height
	return viewport.Sync(b.blogViewport)
}
func renderContent(width int, content string) string {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStyles(glamour.DraculaStyleConfig),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		goreland.LogError("Error initializing glamour")
	}

	str, err := renderer.Render(content)
	if err != nil {
		goreland.LogError("Error rendering content")
		return "Error rendering content!"
	}

	return str

}
func newBlogViewer(width int, height int) *BlogViewer {
	v := viewport.New(width, height)
	return &BlogViewer{
		blogViewport: v,
	}
}
