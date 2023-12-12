package tui

import (
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
func (m Model) selectItem() {
	blog := m.blogList.SelectBlog()
	*m.screenType = blogScreen
	m.blogViewer.SetBlog(m.actualWidth, blog)
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
	case "enter":
		*m.screenType = listScreen
	case "q":
		m.backPressed()
	}
}
func (m Model) handleListKeybindings(msg tea.KeyMsg) {
	switch msg.String() {
	case "q":
		m.backPressed()
	case "enter":
		m.selectItem()
	}
}
func (m Model) handleBlogKeybindings(msg tea.KeyMsg) {
	switch msg.String() {
	case "q":
		m.backPressed()
	}
}
