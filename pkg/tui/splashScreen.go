package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type SplashViewer struct {
	splashContent string
}

func (m Model) viewSplashScreen() string {
	view := strings.Builder{}
	return view.String()
}

func newSplashViewer(width int, height int, splash string) *SplashViewer {

	return &SplashViewer{
		splashContent: splash,
	}
}
func (s *SplashViewer) SetSize(width, height int) tea.Cmd {
	// s.splashViewport.Width = width
	// s.splashViewport.Height = height
	// return viewport.Sync(s.splashViewport)
	return nil

}
func (s *SplashViewer) Update(msg tea.Msg) tea.Cmd {
	return nil
}
