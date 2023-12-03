package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/pspiagicw/goreland"
)
func (m model) viewBlogScreen() string {
	view := strings.Builder{}
	view.WriteString(m.blogViewport.View())
	view.WriteString(fmt.Sprintf("\n%q\n", m.selected))
	return view.String()
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
