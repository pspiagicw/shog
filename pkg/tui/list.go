package tui

import (
	"fmt"
	"strings"
)

func (m model) viewListScreen() string {
	view := strings.Builder{}
	view.WriteString(m.list.View())
	view.WriteString(fmt.Sprintf("\n%q\n", m.selected))
	return view.String()
}
