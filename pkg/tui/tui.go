package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/pspiagicw/shog/pkg/content"
)

const (
	homeScreen = iota
	blogScreen
)

type model struct {
	items        []list.Item
	term         string
	width        int
	height       int
	list         list.Model
	selected     list.Item
	screenType   int
	blogViewport viewport.Model
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.list.SetSize(m.width, m.height)
		m.blogViewport.Height = m.height
		m.blogViewport.Width = m.width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.screenType != homeScreen {
				m.screenType = homeScreen
				return m, nil
			}
			return m, tea.Quit
		case "enter":
			if m.screenType == homeScreen {
				m.selected = m.items[m.list.Index()]
				m.screenType = blogScreen
				m.blogViewport.SetContent(content.GetContentMatter(m.width))
			}
		}
	}
	if m.screenType == homeScreen {
		l, cmd := m.list.Update(msg)
		m.list = l
		return m, cmd
	} else {
		v, cmd := m.blogViewport.Update(msg)
		m.blogViewport = v
		return m, cmd
	}
}

func (m model) View() string {
	view := strings.Builder{}
	switch m.screenType {
	case homeScreen:
		view.WriteString(m.list.View())
		view.WriteString(fmt.Sprintf("\n%q\n", m.selected))
	case blogScreen:
		view.WriteString(m.blogViewport.View())
		view.WriteString(fmt.Sprintf("\n%q\n", m.selected))
	}
	return view.String()
}
func NewModel(pty ssh.Pty) model {
	items := content.GetContent()
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	v := viewport.New(pty.Window.Width, pty.Window.Height)
	m := model{
		items:        items,
		term:         pty.Term,
		width:        pty.Window.Width,
		height:       pty.Window.Height,
		list:         l,
		screenType:   homeScreen,
		blogViewport: v,
	}
	return m

}
func SSHEntry(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, active := s.Pty()
	if !active {
		wish.Fatalln(s, "no active terminal , skipping")
		return nil, nil
	}

	return NewModel(pty), []tea.ProgramOption{tea.WithAltScreen()}
}
