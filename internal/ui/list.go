package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mojoaar/klyde/internal/data"
)

type ListView struct {
	shortcuts []data.Shortcut
	configs   []data.Config
	section   Section
	cursor    int
	focused   bool
	filter    string
	width     int
	height    int
}

func NewListView() ListView {
	return ListView{}
}

func (l *ListView) SetSection(s Section) {
	l.section = s
	l.cursor = 0
}

func (l *ListView) SetData(shortcuts []data.Shortcut, configs []data.Config) {
	l.shortcuts = shortcuts
	l.configs = configs
	// clamp cursor
	if l.cursor >= l.visibleCount() && l.visibleCount() > 0 {
		l.cursor = l.visibleCount() - 1
	}
}

func (l *ListView) SetFilter(f string) {
	l.filter = f
	l.cursor = 0
}

func (l *ListView) SetFocus(f bool) { l.focused = f }
func (l ListView) Focused() bool    { return l.focused }

func (l *ListView) SetSize(w, h int) {
	l.width = w
	l.height = h
}

func (l ListView) visibleCount() int {
	return len(l.filteredShortcuts()) + len(l.filteredConfigs())
}

func (l ListView) filteredShortcuts() []data.Shortcut {
	var out []data.Shortcut
	if l.filter == "" {
		out = make([]data.Shortcut, len(l.shortcuts))
		copy(out, l.shortcuts)
	} else {
		f := strings.ToLower(l.filter)
		for _, s := range l.shortcuts {
			if strings.Contains(strings.ToLower(s.Name), f) ||
				strings.Contains(strings.ToLower(s.Key), f) ||
				strings.Contains(strings.ToLower(s.App), f) {
				out = append(out, s)
			}
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out
}

func (l ListView) filteredConfigs() []data.Config {
	var out []data.Config
	if l.filter == "" {
		out = make([]data.Config, len(l.configs))
		copy(out, l.configs)
	} else {
		f := strings.ToLower(l.filter)
		for _, c := range l.configs {
			if strings.Contains(strings.ToLower(c.Name), f) ||
				strings.Contains(strings.ToLower(c.Path), f) {
				out = append(out, c)
			}
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out
}

// SelectedShortcutID returns the ID of the currently selected shortcut, or "".
func (l ListView) SelectedShortcutID() string {
	items := l.filteredShortcuts()
	if l.section == SectionShortcuts && l.cursor < len(items) {
		return items[l.cursor].ID
	}
	return ""
}

// SelectedConfigID returns the ID of the currently selected config, or "".
func (l ListView) SelectedConfigID() string {
	items := l.filteredConfigs()
	if l.section == SectionConfigs && l.cursor < len(items) {
		return items[l.cursor].ID
	}
	return ""
}

// SelectedShortcut returns the currently selected Shortcut (ok=false if none).
func (l ListView) SelectedShortcut() (data.Shortcut, bool) {
	items := l.filteredShortcuts()
	if l.section == SectionShortcuts && l.cursor < len(items) {
		return items[l.cursor], true
	}
	return data.Shortcut{}, false
}

// SelectedConfig returns the currently selected Config (ok=false if none).
func (l ListView) SelectedConfig() (data.Config, bool) {
	items := l.filteredConfigs()
	if l.section == SectionConfigs && l.cursor < len(items) {
		return items[l.cursor], true
	}
	return data.Config{}, false
}

var listKeys = struct {
	Up   key.Binding
	Down key.Binding
}{
	Up:   key.NewBinding(key.WithKeys("up", "k")),
	Down: key.NewBinding(key.WithKeys("down", "j")),
}

func (l ListView) Update(msg tea.Msg) (ListView, tea.Cmd) {
	if !l.focused {
		return l, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		count := l.visibleCount()
		switch {
		case key.Matches(msg, listKeys.Up):
			if l.cursor > 0 {
				l.cursor--
			}
		case key.Matches(msg, listKeys.Down):
			if l.cursor < count-1 {
				l.cursor++
			}
		}
	}
	return l, nil
}

func (l ListView) View() string {
	switch l.section {
	case SectionShortcuts:
		return l.viewShortcuts()
	case SectionConfigs:
		return l.viewConfigs()
	}
	return ""
}

func (l ListView) viewShortcuts() string {
	items := l.filteredShortcuts()

	// distribute usable width: subtract ListStyle padding (2+2)
	usable := l.width - 4
	if usable < 30 {
		usable = 30
	}
	// Name 40%, Key 35%, App 25%
	nameW := usable * 40 / 100
	keyW := usable * 35 / 100
	appW := usable - nameW - keyW

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		HeaderStyle.Width(nameW).Render("Name"),
		HeaderStyle.Width(keyW).Render("Key"),
		HeaderStyle.Width(appW).Render("App"),
	)
	divider := strings.Repeat("─", nameW+keyW+appW)

	lines := []string{header, HeaderStyle.Render(divider)}

	if len(items) == 0 {
		lines = append(lines, ListItemStyle.Render("  No shortcuts yet. Press 'a' to add one."))
	}

	for i, sc := range items {
		name := truncate(sc.Name, nameW-2)
		k := truncate(sc.Key, keyW-2)
		app := truncate(sc.App, appW-2)
		row := lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Width(nameW).Render(name),
			lipgloss.NewStyle().Width(keyW).Render(k),
			lipgloss.NewStyle().Width(appW).Render(app),
		)
		if i == l.cursor && l.focused {
			lines = append(lines, ListSelectedStyle.Render(row))
		} else {
			lines = append(lines, ListItemStyle.Render(row))
		}
	}

	return ListStyle.Width(l.width).Height(l.height).Render(
		lipgloss.JoinVertical(lipgloss.Left, lines...),
	)
}

func (l ListView) viewConfigs() string {
	items := l.filteredConfigs()

	usable := l.width - 4
	if usable < 30 {
		usable = 30
	}
	// Name 28%, Path 52%, Type 20%
	nameW := usable * 28 / 100
	typeW := usable * 20 / 100
	pathW := usable - nameW - typeW

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		HeaderStyle.Width(nameW).Render("Name"),
		HeaderStyle.Width(pathW).Render("Path"),
		HeaderStyle.Width(typeW).Render("Type"),
	)
	divider := strings.Repeat("─", nameW+pathW+typeW)

	lines := []string{header, HeaderStyle.Render(divider)}

	if len(items) == 0 {
		lines = append(lines, ListItemStyle.Render("  No configs yet. Press 'a' to add one."))
	}

	for i, c := range items {
		name := truncate(c.Name, nameW-2)
		path := truncate(c.Path, pathW-2)
		typ := fmt.Sprintf("%s", c.Type)
		row := lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Width(nameW).Render(name),
			lipgloss.NewStyle().Width(pathW).Render(path),
			lipgloss.NewStyle().Width(typeW).Render(typ),
		)
		if i == l.cursor && l.focused {
			lines = append(lines, ListSelectedStyle.Render(row))
		} else {
			lines = append(lines, ListItemStyle.Render(row))
		}
	}

	return ListStyle.Width(l.width).Height(l.height).Render(
		lipgloss.JoinVertical(lipgloss.Left, lines...),
	)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
