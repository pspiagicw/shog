package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) backPressed() tea.Cmd {
	if *m.screenType == blogScreen {
		*m.screenType = listScreen
		return nil
	} else if *m.screenType == listScreen {
		*m.screenType = splashScreen
		return nil
	} else if *m.screenType == splashScreen {
		return tea.Quit
	}
	return nil
}
func (m Model) selectItem() tea.Cmd {
	if *m.screenType == listScreen {
		blog := m.blogList.SelectBlog()
		*m.screenType = blogScreen
		return m.blogViewer.SetBlog(m.actualWidth, blog)
	}
	return nil
}
func (m Model) handleKeybinding(msg tea.KeyMsg) tea.Cmd {
	if *m.screenType == splashScreen {
		m.handleSplashKeybindings(msg)
	} else if *m.screenType == listScreen {
		m.handleListKeybindings(msg)
	} else {
		m.handleBlogKeybindings(msg)
	}
	return nil
}
func (m Model) handleSplashKeybindings(msg tea.KeyMsg) {
	switch msg.String() {
	case "B":
		fmt.Println("You pressed 'b'")
		*m.screenType = listScreen
	}
}
func (m Model) handleListKeybindings(msg tea.KeyMsg) {
}
func (m Model) handleBlogKeybindings(msg tea.KeyMsg) {
}
