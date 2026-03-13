package app

import (
	"fmt"
	"io/fs"
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

type page struct {
	Title   string
	Summary string
	Content string
}

type model struct {
	width        int
	height       int
	activeTab    int
	projectIndex int
	scroll       int
	content      content
}

func NewModel(contentFS fs.FS) func(ssh.Session) (tea.Model, []tea.ProgramOption) {
	loadedContent := loadContent(contentFS)

	return func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
		m := model{
			content: loadedContent,
		}

		return m, []tea.ProgramOption{
			tea.WithAltScreen(),
		}
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
		case "right", "l", "tab":
			m.activeTab = (m.activeTab + 1) % 3
			m.scroll = 0
		case "left", "h", "shift+tab":
			m.activeTab = (m.activeTab + 2) % 3
			m.scroll = 0
		case "n":
			if m.activeTab == 1 {
				m.projectIndex = (m.projectIndex + 1) % len(m.content.Projects)
				m.scroll = 0
			}
		case "p":
			if m.activeTab == 1 {
				m.projectIndex = (m.projectIndex - 1 + len(m.content.Projects)) % len(m.content.Projects)
				m.scroll = 0
			}
		case "1":
			m.activeTab = 0
			m.scroll = 0
		case "2":
			m.activeTab = 1
			m.scroll = 0
		case "3":
			m.activeTab = 2
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
	current := m.currentPage()
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
			Background(lipgloss.Color("24")).
			Padding(0, 1).
			Bold(true)

	summaryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	bodyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	activeTabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("16")).
			Background(lipgloss.Color("86")).
			Padding(0, 1).
			Bold(true)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Padding(0, 1)
)

func (m model) pageLines(current page, width int) []string {
	lines := []string{
		headerStyle.Render(m.content.Title),
		"",
		m.tabsLine(width),
		"",
	}

	lines = append(lines, wrapStyled(current.Title, width, titleStyle)...)
	lines = append(lines, "")
	lines = append(lines, wrapStyled(current.Summary, width, summaryStyle)...)
	lines = append(lines, "")
	lines = append(lines, renderMarkdown(current.Content, width)...)
	lines = append(lines, "", footerStyle.Render(fmt.Sprintf(
		"%s  •  h/left prev tab  •  l/right next tab  •  j/k scroll  •  q quit",
		m.footerContext(),
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
	return max(0, len(m.pageLines(m.currentPage(), m.contentWidth()))-m.contentHeight())
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

func renderMarkdown(text string, width int) []string {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return wrapStyled(text, width, bodyStyle)
	}

	rendered, err := renderer.Render(strings.TrimSpace(text))
	if err != nil {
		return wrapStyled(text, width, bodyStyle)
	}

	lines := strings.Split(strings.TrimRight(rendered, "\n"), "\n")
	if len(lines) == 0 {
		return []string{""}
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

func (m model) currentPage() page {
	switch m.activeTab {
	case 1:
		return m.content.Projects[m.projectIndex]
	case 2:
		return m.content.Blog
	default:
		return m.content.About
	}
}

func (m model) tabsLine(width int) string {
	tabs := []string{
		m.renderTab("About", 0),
		m.renderTab("Projects", 1),
		m.renderTab("Blog", 2),
	}

	return lipgloss.NewStyle().Width(width).Render(strings.Join(tabs, "  "))
}

func (m model) renderTab(label string, index int) string {
	if m.activeTab == index {
		return activeTabStyle.Render(label)
	}

	return inactiveTabStyle.Render(label)
}

func (m model) footerContext() string {
	if m.activeTab != 1 {
		return fmt.Sprintf("Tab %d/3  •  1/2/3 jump to tab", m.activeTab+1)
	}

	return fmt.Sprintf(
		"Tab %d/3  •  Project %d/%d  •  n/p switch project",
		m.activeTab+1,
		m.projectIndex+1,
		len(m.content.Projects),
	)
}
