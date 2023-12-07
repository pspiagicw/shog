package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

var nameStyle lipgloss.Style = lipgloss.NewStyle().Faint(false).Foreground(lipgloss.Color("#50fa7b")).PaddingLeft(2).PaddingRight(2)
var sloganStyle lipgloss.Style = lipgloss.NewStyle().Faint(true).Italic(true)

func (m Model) View() string {
    view := strings.Builder{}
	switch *m.screenType {
	case listScreen:
		fmt.Println("You are viewing a list screen!")
		view.WriteString(m.viewListScreen())
	case splashScreen:
		fmt.Println("You are viewing a splash screen!")
		return m.viewSplashScreen()
	case blogScreen:
		fmt.Println("You are viewing a blog screen!")
		return m.viewBlogScreen()
	}
    // return displayStyle.Render(view.String())
    return view.String()
}
func (m Model) getHeader() string {
	actualWidth := *m.width - nameStyle.GetPaddingLeft() - nameStyle.GetPaddingRight()
	displayString := nameStyle.Render(m.args.Name) + " - " + sloganStyle.Render(m.args.Slogan)

	return truncate.String(displayString, uint(actualWidth-2))

}
