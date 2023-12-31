package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/pspiagicw/shog/pkg/argparse"
	"github.com/pspiagicw/shog/pkg/content"
)

type ScreenType int

const (
	splashScreen ScreenType = iota
	listScreen
	blogScreen
)

type Model struct {
	width        *int
	height       *int
	actualWidth  int
	actualHeight int
	args         *argparse.Args
	screenType   *ScreenType
	blogList     *BlogList
	splashViewer *SplashViewer
	blogViewer   *BlogViewer
}

func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) SetSize(width, height int) []tea.Cmd {
	cmds := []tea.Cmd{}
	headerHeight := 2
	actualWidth := width - paddingLeft - paddingRight
	actualHeight := height - paddingTop - headerHeight
	m.actualWidth = actualWidth
	m.actualHeight = actualHeight
	cmds = append(cmds, m.blogViewer.SetSize(actualWidth, actualHeight))
	cmds = append(cmds, m.blogList.SetSize(actualWidth, actualHeight))
	cmds = append(cmds, m.splashViewer.SetSize(actualWidth, actualHeight))

	return cmds
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		height, width := msg.Height, msg.Width
		*m.height = height
		*m.width = width
		cmds = append(cmds, m.SetSize(width, height)...)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			m.handleKeybinding(msg)
		}
	}
	if *m.screenType == listScreen {
		cmds = append(cmds, m.blogList.Update(msg))
	} else if *m.screenType == splashScreen {
		cmds = append(cmds, m.splashViewer.Update(msg))
	} else {
		cmds = append(cmds, m.blogViewer.Update(msg))
	}
	return m, tea.Batch(cmds...)
}

func generateItems(blogs []content.Blog) []list.Item {
	items := []list.Item{}
	for _, blog := range blogs {
		items = append(items, blog)
	}
	return items
}
func NewModel(pty ssh.Pty, blogs []content.Blog, splash string, args *argparse.Args) Model {
	width, height := pty.Window.Width, pty.Window.Height
	bl := newBlogList(blogs, width, height)
	bv := newBlogViewer(width, height)
	sv := newSplashViewer(width, height, splash)
	var screen ScreenType = 0
	m := Model{
		width:        &pty.Window.Width,
		height:       &pty.Window.Height,
		blogList:     bl,
		blogViewer:   bv,
		splashViewer: sv,
		screenType:   &screen,
		args:         args,
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

		return NewModel(pty, blogs, splash, args), []tea.ProgramOption{tea.WithAltScreen()}
	}
}
