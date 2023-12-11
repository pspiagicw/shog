package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)


func (m Model) View() string {
    var screen string
	switch *m.screenType {
	case listScreen:
		fmt.Println("You are viewing a list screen!")
		screen = m.viewListScreen()
	case splashScreen:
		fmt.Println("You are viewing a splash screen!")
		return m.viewSplashScreen()
	case blogScreen:
		fmt.Println("You are viewing a blog screen!")
		screen = m.viewBlogScreen()
	}
    screen = lipgloss.JoinVertical(0, m.getHorizontalLine(), screen)
    screen = lipgloss.JoinVertical(0, m.getHeader(), screen)
    return displayStyle.Render(screen)
}
func (m Model) getHeader() string {
	actualWidth := *m.width - nameStyle.GetPaddingLeft() - nameStyle.GetPaddingRight()
	displayString := nameStyle.Render(m.args.Name) + " " + sloganStyle.Render(m.args.Slogan)

	return truncate.String(displayString, uint(actualWidth-2))
}
func (m Model) getHorizontalLine() string {
    var line strings.Builder

    for i := 0; i < *m.width; i++ {
        line.WriteString(dash)
    }
    return line.String()
}
