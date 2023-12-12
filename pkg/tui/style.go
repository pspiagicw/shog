package tui

import "github.com/charmbracelet/lipgloss"

const paddingTop = 2
const paddingLeft = 3
const paddingRight = 3

var displayStyle lipgloss.Style = lipgloss.NewStyle().PaddingLeft(paddingLeft).PaddingTop(paddingTop).PaddingRight(paddingRight)
var nameStyle lipgloss.Style = lipgloss.NewStyle().Faint(false).Foreground(lipgloss.Color("#50fa7b")).PaddingLeft(2).PaddingRight(2)
var sloganStyle lipgloss.Style = lipgloss.NewStyle().Faint(true).Italic(true)

var titleStyle lipgloss.Style = lipgloss.NewStyle().Bold(true)
var infoStyle lipgloss.Style = lipgloss.NewStyle().Faint(true).Italic(true)
var selectedStyle lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#bd93f9"))
var itemStyle lipgloss.Style = lipgloss.NewStyle().BorderLeft(true).BorderStyle(lipgloss.NormalBorder()).PaddingLeft(1)

const dash = "─"
