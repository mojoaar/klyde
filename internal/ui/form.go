package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mojoaar/klyde/internal/data"
)

// FormMode indicates whether the form is for adding or editing.
type FormMode int

const (
	FormModeAdd FormMode = iota
	FormModeEdit
)

// FormSubmitMsg is emitted when the user confirms a form.
type FormSubmitMsg struct {
	Shortcut *data.Shortcut
	Config   *data.Config
}

// FormCancelMsg is emitted when the user cancels the form.
type FormCancelMsg struct{}

type Form struct {
	section  Section
	mode     FormMode
	editID   string
	inputs   []textinput.Model
	focused  int
	width    int
}

var shortcutLabels = []string{"Name", "Key", "App", "Tags"}
var configLabels = []string{"Name", "Path", "Type (file/dir)", "Tags"}

func NewForm(section Section, mode FormMode, width int) Form {
	labels := shortcutLabels
	if section == SectionConfigs {
		labels = configLabels
	}

	inputs := make([]textinput.Model, len(labels))
	for i, label := range labels {
		ti := textinput.New()
		ti.Placeholder = label
		ti.CharLimit = 256
		if i == len(labels)-1 {
			ti.Placeholder = "vim, editor, text"
		}
		if i == 0 {
			ti.Focus()
		}
		inputs[i] = ti
	}

	return Form{
		section: section,
		mode:    mode,
		inputs:  inputs,
		focused: 0,
		width:   width,
	}
}

// Prefill populates the form from an existing Shortcut for editing.
func (f *Form) PrefillShortcut(sc data.Shortcut) {
	f.editID = sc.ID
	f.inputs[0].SetValue(sc.Name)
	f.inputs[1].SetValue(sc.Key)
	f.inputs[2].SetValue(sc.App)
	f.inputs[3].SetValue(strings.Join(sc.Tags, ", "))
}

// Prefill populates the form from an existing Config for editing.
func (f *Form) PrefillConfig(c data.Config) {
	f.editID = c.ID
	f.inputs[0].SetValue(c.Name)
	f.inputs[1].SetValue(c.Path)
	f.inputs[2].SetValue(string(c.Type))
	f.inputs[3].SetValue(strings.Join(c.Tags, ", "))
}

func (f Form) Update(msg tea.Msg) (Form, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return f, func() tea.Msg { return FormCancelMsg{} }
		case "enter":
			if f.focused == len(f.inputs)-1 {
				return f, f.submit()
			}
			f.inputs[f.focused].Blur()
			f.focused++
			f.inputs[f.focused].Focus()
			return f, nil
		case "tab", "down":
			f.inputs[f.focused].Blur()
			f.focused = (f.focused + 1) % len(f.inputs)
			f.inputs[f.focused].Focus()
			return f, nil
		case "shift+tab", "up":
			f.inputs[f.focused].Blur()
			f.focused = (f.focused - 1 + len(f.inputs)) % len(f.inputs)
			f.inputs[f.focused].Focus()
			return f, nil
		case "ctrl+s":
			return f, f.submit()
		}
	}

	var cmds []tea.Cmd
	for i := range f.inputs {
		var cmd tea.Cmd
		f.inputs[i], cmd = f.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}
	return f, tea.Batch(cmds...)
}

func (f Form) submit() tea.Cmd {
	return func() tea.Msg {
		tags := parseTags(f.inputs[3].Value())
		if f.section == SectionShortcuts {
			sc := data.Shortcut{
				ID:   f.editID,
				Name: f.inputs[0].Value(),
				Key:  f.inputs[1].Value(),
				App:  f.inputs[2].Value(),
				Tags: tags,
			}
			return FormSubmitMsg{Shortcut: &sc}
		}
		entryType := data.EntryTypeFile
		if strings.ToLower(strings.TrimSpace(f.inputs[2].Value())) == "directory" ||
			strings.ToLower(strings.TrimSpace(f.inputs[2].Value())) == "dir" {
			entryType = data.EntryTypeDirectory
		}
		c := data.Config{
			ID:   f.editID,
			Name: f.inputs[0].Value(),
			Path: f.inputs[1].Value(),
			Type: entryType,
			Tags: tags,
		}
		return FormSubmitMsg{Config: &c}
	}
}

func (f Form) View() string {
	labels := shortcutLabels
	if f.section == SectionConfigs {
		labels = configLabels
	}

	title := "Add"
	if f.mode == FormModeEdit {
		title = "Edit"
	}
	sectionName := "Shortcut"
	if f.section == SectionConfigs {
		sectionName = "Config"
	}

	var rows []string
	rows = append(rows, FormTitleStyle.Render(title+" "+sectionName))

	inputWidth := f.width - 12
	if inputWidth < 20 {
		inputWidth = 20
	}

	for i, ti := range f.inputs {
		label := FormLabelStyle.Render(labels[i])
		var inputView string
		if i == f.focused {
			inputView = FormActiveInputStyle.Width(inputWidth).Render(ti.View())
		} else {
			inputView = FormInputStyle.Width(inputWidth).Render(ti.View())
		}
		rows = append(rows, lipgloss.JoinVertical(lipgloss.Left, label, inputView))
		rows = append(rows, "") // spacing between fields
	}

	rows = append(rows, "")
	rows = append(rows, HelpDescStyle.Render("tab/↑↓ navigate  •  enter next  •  ctrl+s submit  •  esc cancel"))

	return lipgloss.NewStyle().
		Padding(2, 4).
		Width(f.width).
		Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}

func parseTags(s string) []string {
	if strings.TrimSpace(s) == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	var tags []string
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}
