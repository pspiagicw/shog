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
	Meta    string
}

type model struct {
	width        int
	height       int
	activeTab    int
	projectIndex int
	blogIndex    int
	projectList  bool
	blogList     bool
	showHelp     bool
	scroll       int
	content      content
}

func NewModel(contentFS fs.FS) func(ssh.Session) (tea.Model, []tea.ProgramOption) {
	loadedContent := loadContent(contentFS)

	return func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
		m := model{
			content:     loadedContent,
			projectList: true,
			blogList:    true,
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
		case "?":
			m.showHelp = !m.showHelp
			m.scroll = 0
			return m, nil
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		if m.showHelp {
			return m.updateScrollableKey(msg.String()), nil
		}

		switch msg.String() {
		case "right", "l", "tab":
			m.activeTab = (m.activeTab + 1) % 3
			m.projectList = true
			m.blogList = true
			m.scroll = 0
		case "left", "h", "shift+tab":
			m.activeTab = (m.activeTab + 2) % 3
			m.projectList = true
			m.blogList = true
			m.scroll = 0
		case "1":
			m.activeTab = 0
			m.projectList = true
			m.blogList = true
			m.scroll = 0
		case "2":
			m.activeTab = 1
			m.projectList = true
			m.scroll = 0
		case "3":
			m.activeTab = 2
			m.blogList = true
			m.scroll = 0
		case "enter":
			if m.activeTab == 1 {
				m.projectList = false
				m.scroll = 0
			} else if m.activeTab == 2 {
				m.blogList = false
				m.scroll = 0
			}
		case "esc", "backspace":
			if m.activeTab == 1 && !m.projectList {
				m.projectList = true
				m.scroll = 0
			} else if m.activeTab == 2 && !m.blogList {
				m.blogList = true
				m.scroll = 0
			}
		case "n":
			if m.activeTab == 1 {
				m.projectIndex = (m.projectIndex + 1) % len(m.content.Projects)
				m.scroll = 0
			} else if m.activeTab == 2 {
				m.blogIndex = (m.blogIndex + 1) % len(m.content.Blogs)
				m.scroll = 0
			}
		case "p":
			if m.activeTab == 1 {
				m.projectIndex = (m.projectIndex - 1 + len(m.content.Projects)) % len(m.content.Projects)
				m.scroll = 0
			} else if m.activeTab == 2 {
				m.blogIndex = (m.blogIndex - 1 + len(m.content.Blogs)) % len(m.content.Blogs)
				m.scroll = 0
			}
		case "down", "j":
			if m.activeTab == 1 && m.projectList {
				m.projectIndex = min(m.projectIndex+1, len(m.content.Projects)-1)
				m.scroll = 0
			} else if m.activeTab == 2 && m.blogList {
				m.blogIndex = min(m.blogIndex+1, len(m.content.Blogs)-1)
				m.scroll = 0
			} else {
				m = m.updateScrollableKey(msg.String())
			}
		case "up", "k":
			if m.activeTab == 1 && m.projectList {
				m.projectIndex = max(m.projectIndex-1, 0)
				m.scroll = 0
			} else if m.activeTab == 2 && m.blogList {
				m.blogIndex = max(m.blogIndex-1, 0)
				m.scroll = 0
			} else {
				m = m.updateScrollableKey(msg.String())
			}
		default:
			m = m.updateScrollableKey(msg.String())
		}
	}

	return m, nil
}

func (m model) updateScrollableKey(key string) model {
	switch key {
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

	return m
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

	metaStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")).
			Italic(true)

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
	if current.Summary != "" {
		lines = append(lines, wrapStyled(current.Summary, width, summaryStyle)...)
		lines = append(lines, "")
	}
	if current.Meta != "" {
		lines = append(lines, wrapStyled(current.Meta, width, metaStyle)...)
		lines = append(lines, "")
	}
	lines = append(lines, renderMarkdown(current.Content, width)...)
	lines = append(lines, "", footerStyle.Render(fmt.Sprintf(
		"%s  •  h/left prev tab  •  l/right next tab  •  ? help  •  q quit",
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

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func (m model) currentPage() page {
	if m.showHelp {
		return m.helpPage()
	}

	switch m.activeTab {
	case 1:
		if m.projectList {
			return m.projectIndexPage()
		}
		return m.content.Projects[m.projectIndex]
	case 2:
		if m.blogList {
			return m.blogIndexPage()
		}
		return m.content.Blogs[m.blogIndex]
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
	if m.showHelp {
		return "Help  •  esc or ? close"
	}

	if m.activeTab == 2 {
		if m.blogList {
			return fmt.Sprintf("Tab %d/3  •  %d posts  •  j/k move  •  enter open", m.activeTab+1, len(m.content.Blogs))
		}

		return fmt.Sprintf("Tab %d/3  •  Post %d/%d  •  n/p switch post  •  esc back", m.activeTab+1, m.blogIndex+1, len(m.content.Blogs))
	}

	if m.activeTab != 1 {
		return fmt.Sprintf("Tab %d/3  •  1/2/3 jump to tab", m.activeTab+1)
	}

	if m.projectList {
		return fmt.Sprintf("Tab %d/3  •  %d projects  •  j/k move  •  enter open", m.activeTab+1, len(m.content.Projects))
	}

	return fmt.Sprintf("Tab %d/3  •  Project %d/%d  •  n/p switch project  •  esc back", m.activeTab+1, m.projectIndex+1, len(m.content.Projects))
}

func (m model) projectIndexPage() page {
	items := make([]string, 0, len(m.content.Projects)*2)
	for i, project := range m.content.Projects {
		prefix := "  "
		if i == m.projectIndex {
			prefix = "> "
		}

		items = append(items, fmt.Sprintf("%s%s", prefix, project.Title))
		items = append(items, fmt.Sprintf("    %s", firstLine(project.Summary)))
	}

	return page{
		Title:   "Projects",
		Summary: "Project index",
		Content: strings.TrimSpace(strings.Join(items, "\n")),
	}
}

func (m model) blogIndexPage() page {
	items := make([]string, 0, len(m.content.Blogs)*3)
	for i, post := range m.content.Blogs {
		prefix := "  "
		if i == m.blogIndex {
			prefix = "> "
		}

		line := fmt.Sprintf("%s%s", prefix, post.Title)
		if post.Meta != "" {
			line = fmt.Sprintf("%s (%s)", line, post.Meta)
		}
		items = append(items, line)
		if post.Summary != "" {
			items = append(items, fmt.Sprintf("    %s", post.Summary))
		}
		items = append(items, "")
	}

	return page{
		Title:   "Blog",
		Summary: "Blog index",
		Content: strings.TrimSpace(strings.Join(items, "\n")),
	}
}

func (m model) helpPage() page {
	return page{
		Title:   "Help",
		Summary: "Keyboard shortcuts",
		Content: strings.TrimSpace(strings.Join([]string{
			"## Global",
			"",
			"- `1`, `2`, `3`: jump to About, Projects, Blog",
			"- `h` / `left`, `l` / `right`: switch tabs",
			"- `?`: toggle help",
			"- `q`: quit",
			"",
			"## Index views",
			"",
			"- `j` / `k`: move selection",
			"- `enter`: open selected item",
			"",
			"## Detail views",
			"",
			"- `n` / `p`: next or previous project/post",
			"- `j` / `k`: scroll",
			"- `space` / `f`: page down",
			"- `b`: page up",
			"- `g` / `G`: top / bottom",
			"- `esc` / `backspace`: return to index",
		}, "\n")),
	}
}

func firstLine(value string) string {
	lines := strings.Split(strings.TrimSpace(value), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return ""
	}

	return lines[0]
}
