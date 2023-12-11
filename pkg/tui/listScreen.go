package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pspiagicw/shog/pkg/content"
)

type BlogList struct {
	blogList list.Model
	selected content.Blog
	blogs    []content.Blog
}

func (b *BlogList) GetIndex() int {
	return b.blogList.Index()
}

func (b *BlogList) SelectBlog() content.Blog {
	b.selected = b.blogs[b.GetIndex()]
	return b.selected
}
func (b *BlogList) Update(msg tea.Msg) tea.Cmd {
	l, cmd := b.blogList.Update(msg)
	b.blogList = l
	return cmd
}
func (m Model) viewListScreen() string {
	view := strings.Builder{}
	view.WriteString(m.blogList.blogList.View())
	return view.String()
}
func (b *BlogList) SetSize(width, height int) tea.Cmd {
	b.blogList.SetSize(width, height)
	return nil
}

type BlogDelegate struct {
}

func (b BlogDelegate) Height() int {
	return 2
}

func (b BlogDelegate) Spacing() int {
	return 2
}

func (b BlogDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
func (b BlogDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	blog := item.(content.Blog)
	s := strings.Builder{}
	s.WriteString(blog.BlogTitle)
	s.WriteString(blog.BlogTitle)
	_, _ = fmt.Fprintf(w, s.String())
}
func newBlogList(blogs []content.Blog, width, height int) *BlogList {
	items := generateItems(blogs)
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	return &BlogList{
		blogList: l,
		blogs:    blogs,
	}
}
