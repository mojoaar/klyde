package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmYesMsg is emitted when the user confirms deletion.
type ConfirmYesMsg struct{ ID string }

// ConfirmNoMsg is emitted when the user cancels.
type ConfirmNoMsg struct{}

type Confirm struct {
	id      string
	label   string
	choice  bool // true = yes focused
	visible bool
}

func NewConfirm() Confirm { return Confirm{} }

func (c *Confirm) Show(id, label string) {
	c.id = id
	c.label = label
	c.choice = false
	c.visible = true
}

func (c *Confirm) Hide() { c.visible = false }
func (c Confirm) Visible() bool { return c.visible }

var confirmKeys = struct {
	Left  key.Binding
	Right key.Binding
	Yes   key.Binding
	No    key.Binding
}{
	Left:  key.NewBinding(key.WithKeys("left", "h")),
	Right: key.NewBinding(key.WithKeys("right", "l")),
	Yes:   key.NewBinding(key.WithKeys("y", "Y")),
	No:    key.NewBinding(key.WithKeys("n", "N", "esc")),
}

func (c Confirm) Update(msg tea.Msg) (Confirm, tea.Cmd) {
	if !c.visible {
		return c, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, confirmKeys.Left):
			c.choice = true
		case key.Matches(msg, confirmKeys.Right):
			c.choice = false
		case key.Matches(msg, confirmKeys.Yes):
			c.visible = false
			id := c.id
			return c, func() tea.Msg { return ConfirmYesMsg{ID: id} }
		case key.Matches(msg, confirmKeys.No):
			c.visible = false
			return c, func() tea.Msg { return ConfirmNoMsg{} }
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			c.visible = false
			if c.choice {
				id := c.id
				return c, func() tea.Msg { return ConfirmYesMsg{ID: id} }
			}
			return c, func() tea.Msg { return ConfirmNoMsg{} }
		}
	}
	return c, nil
}

func (c Confirm) View(totalWidth, totalHeight int) string {
	yes := ConfirmNoStyle.Render("[ Yes ]")
	no := ConfirmNoStyle.Render("[ No ]")
	if c.choice {
		yes = ConfirmYesStyle.Render("[ Yes ]")
	} else {
		no = lipgloss.NewStyle().
			Background(lipgloss.Color("#374151")).
			Foreground(lipgloss.Color("#F9FAFB")).
			Bold(true).
			Padding(0, 2).
			Render("[ No ]")
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, yes, "  ", no)
	content := lipgloss.JoinVertical(lipgloss.Center,
		ConfirmTitleStyle.Render("Delete"),
		HelpDescStyle.Render(c.label),
		"",
		buttons,
		"",
		HelpDescStyle.Render("←/→ select  •  enter confirm  •  esc cancel"),
	)

	box := ConfirmBoxStyle.Width(44).Render(content)

	return lipgloss.Place(totalWidth, totalHeight, lipgloss.Center, lipgloss.Center, box)
}
