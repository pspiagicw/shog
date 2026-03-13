package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

type page struct {
	Title   string
	Summary string
	Body    []string
}

type model struct {
	width  int
	height int
	index  int
	pages  []page
}

func NewModel(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	m := model{
		pages: seedPages(),
	}

	return m, []tea.ProgramOption{
		tea.WithAltScreen(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n":
			m.index = (m.index + 1) % len(m.pages)
		case "left", "h", "p":
			m.index = (m.index - 1 + len(m.pages)) % len(m.pages)
		}
	}

	return m, nil
}

func (m model) View() string {
	if len(m.pages) == 0 {
		return "No pages configured."
	}

	current := m.pages[m.index]
	header := headerStyle.Render("shog")
	title := titleStyle.Render(current.Title)
	summary := summaryStyle.Render(current.Summary)

	body := make([]string, 0, len(current.Body))
	for _, line := range current.Body {
		body = append(body, bodyStyle.Render(line))
	}

	footer := footerStyle.Render(fmt.Sprintf("Page %d/%d  •  h/left prev  •  l/right next  •  q quit", m.index+1, len(m.pages)))
	content := strings.Join([]string{
		header,
		title,
		summary,
		strings.Join(body, "\n"),
		footer,
	}, "\n\n")

	return frameStyle.Width(max(60, m.width-4)).Render(content)
}

func seedPages() []page {
	return []page{
		{
			Title:   "About",
			Summary: "A short intro delivered over SSH.",
			Body: []string{
				"Hi, I'm Your Name.",
				"I build software, explore ideas, and use this terminal app as a compact portfolio.",
				"Use the arrow keys or h/l to move between pages.",
			},
		},
		{
			Title:   "Project One",
			Summary: "Replace this with a real flagship project.",
			Body: []string{
				"Problem: Describe the problem space in one sentence.",
				"Build: Note the stack, architecture, or unusual constraint.",
				"Result: Summarize the outcome, adoption, or what you learned.",
			},
		},
		{
			Title:   "Project Two",
			Summary: "A second page keeps the navigation model honest.",
			Body: []string{
				"Highlight a contrasting project here.",
				"This can later become a richer card, timeline, or interactive selector.",
				"A future blog page can slot into the same page list without changing the app shape.",
			},
		},
	}
}

var (
	frameStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63"))

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Bold(true)

	summaryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	bodyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))
)

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
