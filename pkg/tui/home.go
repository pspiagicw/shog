package tui

import "strings"

func (m model) viewHomeScreen() string {
	view := strings.Builder{}
    view.WriteString(m.homeViewport.View())
    return view.String()
}
