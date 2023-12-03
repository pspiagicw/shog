package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/pspiagicw/shog/pkg/argparse"
	"github.com/pspiagicw/shog/pkg/content"
)

type ScreenType int

const (
	homeScreen ScreenType = iota
	listScreen
	blogScreen
)

type model struct {
	width         int
	height        int
	list          list.Model
	selected      content.Blog
	screenType    ScreenType
	blogViewport  viewport.Model
	blogs         []content.Blog
	args          *argparse.Args
	homeViewport  viewport.Model
	splashContent string
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
		m.homeViewport.Height = m.height
		m.homeViewport.Width = m.width
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
				m.selected = m.blogs[m.list.Index()]
				m.screenType = blogScreen
				content := m.blogs[m.list.Index()].Content
				m.blogViewport.SetContent(renderContent(m.width, content))

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
	switch m.screenType {
	case homeScreen:
		return m.viewHomeScreen()
	case listScreen:
		return m.viewListScreen()
	case blogScreen:
		return m.viewBlogScreen()
	}
	return ""
}
func generateItems(blogs []content.Blog) []list.Item {
	items := []list.Item{}
	for _, blog := range blogs {
		items = append(items, blog)
	}
	return items
}
func NewModel(pty ssh.Pty, blogs []content.Blog, splash string) model {
	items := generateItems(blogs)
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	v := viewport.New(pty.Window.Width, pty.Window.Height)
	h := viewport.New(pty.Window.Width, pty.Window.Height)
	h.SetContent(splash)
	m := model{
		width:         pty.Window.Width,
		height:        pty.Window.Height,
		list:          l,
		screenType:    homeScreen,
		blogViewport:  v,
		blogs:         blogs,
		homeViewport:  h,
		splashContent: splash,
	}
	return m

}
func EntryGenerator(args *argparse.Args, blogs []content.Blog, splash string) func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		pty, _, active := s.Pty()
		if !active {
			wish.Fatalln(s, "no active terminal , skipping")
			return nil, nil
		}

		return NewModel(pty, blogs, splash), []tea.ProgramOption{tea.WithAltScreen()}
	}
}
