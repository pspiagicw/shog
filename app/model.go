package app

import (
	"fmt"
	"strings"
	"unicode/utf8"

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
	scroll int
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
		m.clampScroll()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n":
			m.index = (m.index + 1) % len(m.pages)
			m.scroll = 0
		case "left", "h", "p":
			m.index = (m.index - 1 + len(m.pages)) % len(m.pages)
			m.scroll = 0
		case "down", "j":
			m.scroll++
			m.clampScroll()
		case "up", "k":
			m.scroll--
			m.clampScroll()
		case "pgdown", "f", " ":
			m.scroll += max(1, m.contentHeight()-1)
			m.clampScroll()
		case "pgup", "b":
			m.scroll -= max(1, m.contentHeight()-1)
			m.clampScroll()
		case "home", "g":
			m.scroll = 0
		case "end", "G":
			m.scroll = m.maxScroll()
		}
	}

	return m, nil
}

func (m model) View() string {
	if len(m.pages) == 0 {
		return "No pages configured."
	}

	current := m.pages[m.index]
	contentWidth := m.contentWidth()
	lines := m.pageLines(current, contentWidth)
	start := min(m.scroll, max(0, len(lines)-m.contentHeight()))
	end := min(len(lines), start+m.contentHeight())
	visible := append([]string{}, lines[start:end]...)

	for len(visible) < m.contentHeight() {
		visible = append(visible, "")
	}

	return frameStyle.
		Width(m.frameWidth()).
		Height(m.frameHeight()).
		Render(strings.Join(visible, "\n"))
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

func (m model) pageLines(current page, width int) []string {
	lines := []string{
		headerStyle.Render("shog"),
		"",
	}

	lines = append(lines, wrapStyled(current.Title, width, titleStyle)...)
	lines = append(lines, "")
	lines = append(lines, wrapStyled(current.Summary, width, summaryStyle)...)
	lines = append(lines, "")

	for i, line := range current.Body {
		lines = append(lines, wrapStyled(line, width, bodyStyle)...)
		if i < len(current.Body)-1 {
			lines = append(lines, "")
		}
	}

	lines = append(lines, "", footerStyle.Render(fmt.Sprintf(
		"Page %d/%d  •  h/left prev  •  l/right next  •  j/k scroll  •  q quit",
		m.index+1,
		len(m.pages),
	)))

	return lines
}

func (m model) frameWidth() int {
	if m.width <= 0 {
		return 76
	}

	return max(60, m.width-4)
}

func (m model) frameHeight() int {
	if m.height <= 0 {
		return 20
	}

	return max(12, m.height-2)
}

func (m model) contentWidth() int {
	return max(20, m.frameWidth()-frameStyle.GetHorizontalFrameSize())
}

func (m model) contentHeight() int {
	return max(5, m.frameHeight()-frameStyle.GetVerticalFrameSize())
}

func (m model) maxScroll() int {
	if len(m.pages) == 0 {
		return 0
	}

	return max(0, len(m.pageLines(m.pages[m.index], m.contentWidth()))-m.contentHeight())
}

func (m *model) clampScroll() {
	m.scroll = max(0, min(m.scroll, m.maxScroll()))
}

func wrapStyled(text string, width int, style lipgloss.Style) []string {
	plain := wrapText(text, width)
	lines := make([]string, 0, len(plain))
	for _, line := range plain {
		lines = append(lines, style.Width(width).Render(line))
	}

	return lines
}

func wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	paragraphs := strings.Split(text, "\n")
	lines := make([]string, 0, len(paragraphs))

	for _, paragraph := range paragraphs {
		if paragraph == "" {
			lines = append(lines, "")
			continue
		}

		words := strings.Fields(paragraph)
		if len(words) == 0 {
			lines = append(lines, "")
			continue
		}

		current := words[0]
		for _, word := range words[1:] {
			candidate := current + " " + word
			if utf8.RuneCountInString(candidate) <= width {
				current = candidate
				continue
			}

			lines = append(lines, current)
			current = word

			for utf8.RuneCountInString(current) > width {
				chunk, rest := splitAtWidth(current, width)
				lines = append(lines, chunk)
				current = rest
			}
		}

		for utf8.RuneCountInString(current) > width {
			chunk, rest := splitAtWidth(current, width)
			lines = append(lines, chunk)
			current = rest
		}

		lines = append(lines, current)
	}

	return lines
}

func splitAtWidth(text string, width int) (string, string) {
	if utf8.RuneCountInString(text) <= width {
		return text, ""
	}

	count := 0
	for i := range text {
		if count == width {
			return text[:i], text[i:]
		}
		count++
	}

	return text, ""
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
