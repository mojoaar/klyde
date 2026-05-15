package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Section represents a top-level navigation item in the sidebar.
type Section int

const (
	SectionShortcuts Section = iota
	SectionConfigs
)

var sectionLabels = []string{"Shortcuts", "Configs"}

type Sidebar struct {
	Current Section
	focused bool
}

func NewSidebar() Sidebar {
	return Sidebar{Current: SectionShortcuts}
}

func (s *Sidebar) SetFocus(f bool) { s.focused = f }
func (s Sidebar) Focused() bool    { return s.focused }

var sidebarKeys = struct {
	Up   key.Binding
	Down key.Binding
}{
	Up:   key.NewBinding(key.WithKeys("up", "k")),
	Down: key.NewBinding(key.WithKeys("down", "j")),
}

func (s Sidebar) Update(msg tea.Msg) (Sidebar, tea.Cmd) {
	if !s.focused {
		return s, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, sidebarKeys.Up):
			if int(s.Current) > 0 {
				s.Current--
			}
		case key.Matches(msg, sidebarKeys.Down):
			if int(s.Current) < len(sectionLabels)-1 {
				s.Current++
			}
		}
	}
	return s, nil
}

func (s Sidebar) View(height int) string {
	var items []string
	for i, label := range sectionLabels {
		var item string
		if Section(i) == s.Current {
			item = SidebarActiveStyle.Render("▶ " + label)
		} else {
			item = SidebarItemStyle.Render("  " + label)
		}
		items = append(items, item)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)

	style := SidebarStyle
	if s.focused {
		style = SidebarFocusBorderStyle
	}
	return style.Height(height).Render(content)
}
