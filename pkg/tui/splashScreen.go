package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pspiagicw/shog/pkg/content"
)

type SplashViewer struct {
	splashContent string
	height        int
	width         int
}

func (m Model) viewSplashScreen() string {
	view := strings.Builder{}
	view.WriteString(renderContent(m.actualWidth, content.DEFAULT_SPLASH))
	return view.String()
}

func newSplashViewer(width int, height int, splash string) *SplashViewer {

	return &SplashViewer{
		splashContent: splash,
	}
}
func (s *SplashViewer) SetSize(width, height int) tea.Cmd {
	s.height = height
	s.width = width
	return nil

}
func (s *SplashViewer) Update(msg tea.Msg) tea.Cmd {
	return nil
}
