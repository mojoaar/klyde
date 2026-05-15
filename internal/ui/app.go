package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mojoaar/klyde/internal/data"
)

// appMode tracks what the user is currently doing.
type appMode int

const (
	modeNav     appMode = iota // browsing sidebar + list
	modeForm                   // add/edit form overlay
	modeConfirm                // delete confirmation overlay
	modeSearch                 // typing a search filter
	modeHelp                   // help overlay
)

type App struct {
	store   *data.Store
	version string
	mode    appMode
	sidebar Sidebar
	list    ListView
	form    Form
	confirm Confirm
	search  textinput.Model

	focus     focusTarget
	width     int
	height    int
	statusMsg string
}

type focusTarget int

const (
	focusSidebar focusTarget = iota
	focusList
)

func NewApp(store *data.Store, version string) App {
	sb := NewSidebar()
	lv := NewListView()
	si := textinput.New()
	si.Placeholder = "search…"
	si.CharLimit = 128

	app := App{
		store:   store,
		version: version,
		sidebar: sb,
		list:    lv,
		confirm: NewConfirm(),
		search:  si,
		focus:   focusList,
	}
	app.syncSection()
	return app
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.list.SetSize(a.width-22, a.height-bannerLines-1)
		return a, nil

	// — Form events —
	case FormSubmitMsg:
		if msg.Shortcut != nil {
			if msg.Shortcut.ID == "" {
				a.store.AddShortcut(*msg.Shortcut)
			} else {
				a.store.UpdateShortcut(*msg.Shortcut)
			}
		} else if msg.Config != nil {
			if msg.Config.ID == "" {
				a.store.AddConfig(*msg.Config)
			} else {
				a.store.UpdateConfig(*msg.Config)
			}
		}
		_ = a.store.Save()
		a.list.SetSection(a.sidebar.Current)
		a.list.SetData(a.store.Shortcuts, a.store.Configs)
		a.mode = modeNav
		a.sidebar.SetFocus(a.focus == focusSidebar)
		a.list.SetFocus(a.focus == focusList)
		return a, nil

	case FormCancelMsg:
		a.mode = modeNav
		a.sidebar.SetFocus(a.focus == focusSidebar)
		a.list.SetFocus(a.focus == focusList)
		return a, nil

	// — Confirm events —
	case ConfirmYesMsg:
		if a.sidebar.Current == SectionShortcuts {
			a.store.DeleteShortcut(msg.ID)
		} else {
			a.store.DeleteConfig(msg.ID)
		}
		_ = a.store.Save()
		a.list.SetSection(a.sidebar.Current)
		a.list.SetData(a.store.Shortcuts, a.store.Configs)
		a.mode = modeNav
		a.sidebar.SetFocus(a.focus == focusSidebar)
		a.list.SetFocus(a.focus == focusList)
		return a, nil

	case ConfirmNoMsg:
		a.mode = modeNav
		a.sidebar.SetFocus(a.focus == focusSidebar)
		a.list.SetFocus(a.focus == focusList)
		return a, nil
	}

	// Delegate key events to the active mode.
	switch a.mode {
	case modeForm:
		return a.updateForm(msg)
	case modeConfirm:
		return a.updateConfirm(msg)
	case modeSearch:
		return a.updateSearch(msg)
	case modeHelp:
		return a.updateHelp(msg)
	default:
		return a.updateNav(msg)
	}
}

// — Mode-specific updaters —

func (a App) updateNav(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		case "?":
			a.mode = modeHelp
			return a, nil
		case "tab":
			if a.focus == focusSidebar {
				a.focus = focusList
			} else {
				a.focus = focusSidebar
			}
			a.syncFocus()
			return a, nil
		case "shift+tab":
			if a.focus == focusList {
				a.focus = focusSidebar
			} else {
				a.focus = focusList
			}
			a.syncFocus()
			return a, nil
		case "a":
			a.openForm(FormModeAdd)
			return a, nil
		case "e":
			a.openForm(FormModeEdit)
			return a, nil
		case "d":
			a.openConfirm()
			return a, nil
		case "/":
			a.mode = modeSearch
			a.search.SetValue("")
			a.search.Focus()
			return a, textinput.Blink
		case "esc":
			// clear search
			a.list.SetFilter("")
			a.search.SetValue("")
			return a, nil
		}

		// When focus is on sidebar, update it; section change syncs list.
		if a.focus == focusSidebar {
			var cmd tea.Cmd
			prev := a.sidebar.Current
			a.sidebar, cmd = a.sidebar.Update(msg)
			if a.sidebar.Current != prev {
				a.syncSection()
			}
			return a, cmd
		}

		// Otherwise update list.
		var cmd tea.Cmd
		a.list, cmd = a.list.Update(msg)
		return a, cmd
	}
	return a, nil
}

func (a App) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.form, cmd = a.form.Update(msg)
	return a, cmd
}

func (a App) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.confirm, cmd = a.confirm.Update(msg)
	return a, cmd
}

func (a App) updateSearch(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "esc":
			a.mode = modeNav
			a.search.Blur()
			return a, nil
		}
	}
	var cmd tea.Cmd
	a.search, cmd = a.search.Update(msg)
	a.list.SetFilter(a.search.Value())
	return a, cmd
}

func (a App) updateHelp(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "?", "esc", "q":
			a.mode = modeNav
		}
	}
	return a, nil
}

// — Helpers —

func (a *App) syncSection() {
	a.list.SetSection(a.sidebar.Current)
	a.list.SetData(a.store.Shortcuts, a.store.Configs)
}

func (a *App) syncFocus() {
	a.sidebar.SetFocus(a.focus == focusSidebar)
	a.list.SetFocus(a.focus == focusList)
}

func (a *App) openForm(mode FormMode) {
	f := NewForm(a.sidebar.Current, mode, a.width)
	if mode == FormModeEdit {
		if a.sidebar.Current == SectionShortcuts {
			sc, ok := a.list.SelectedShortcut()
			if !ok {
				return
			}
			f.PrefillShortcut(sc)
		} else {
			c, ok := a.list.SelectedConfig()
			if !ok {
				return
			}
			f.PrefillConfig(c)
		}
	}
	a.form = f
	a.mode = modeForm
	a.sidebar.SetFocus(false)
	a.list.SetFocus(false)
}

func (a *App) openConfirm() {
	if a.sidebar.Current == SectionShortcuts {
		sc, ok := a.list.SelectedShortcut()
		if !ok {
			return
		}
		a.confirm.Show(sc.ID, fmt.Sprintf(`"%s"`, sc.Name))
	} else {
		c, ok := a.list.SelectedConfig()
		if !ok {
			return
		}
		a.confirm.Show(c.ID, fmt.Sprintf(`"%s"`, c.Name))
	}
	a.mode = modeConfirm
	a.sidebar.SetFocus(false)
	a.list.SetFocus(false)
}

func (a *App) handleFormSubmit(msg FormSubmitMsg) {
	if msg.Shortcut != nil {
		if msg.Shortcut.ID == "" {
			a.store.AddShortcut(*msg.Shortcut)
		} else {
			a.store.UpdateShortcut(*msg.Shortcut)
		}
	} else if msg.Config != nil {
		if msg.Config.ID == "" {
			a.store.AddConfig(*msg.Config)
		} else {
			a.store.UpdateConfig(*msg.Config)
		}
	}
	_ = a.store.Save()
	a.syncSection()
}

func (a *App) handleDelete(id string) {
	if a.sidebar.Current == SectionShortcuts {
		a.store.DeleteShortcut(id)
	} else {
		a.store.DeleteConfig(id)
	}
	_ = a.store.Save()
	a.syncSection()
}

// — View —

func (a App) View() string {
	if a.width == 0 {
		return "Loading…"
	}

	switch a.mode {
	case modeHelp:
		return a.viewHelp()
	case modeForm:
		return a.viewForm()
	case modeConfirm:
		return a.viewConfirmOverlay()
	}

	return a.viewMain()
}

// bannerLines is the fixed height of the banner: 1 padding top + 5 art lines + 1 padding bottom.
const bannerLines = 7

func (a App) viewMain() string {
	// ensure focus state is reflected
	sbFocused := a.focus == focusSidebar
	a.sidebar.SetFocus(sbFocused)
	a.list.SetFocus(!sbFocused)

	contentHeight := a.height - bannerLines - 1 // 1 = help bar

	sidebar := a.sidebar.View(contentHeight)
	list := a.list.View()

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, list)

	helpBar := a.viewHelpBar()

	return lipgloss.JoinVertical(lipgloss.Left, a.viewBanner(), body, helpBar)
}

func (a App) viewBanner() string {
lines := []string{
"   __    __        __   ",
"  / /__ / /_ _____/ /__ ",
" /  '_// / // / _  / -_)",
`/_/\_\/_/\_, /\_,_/\__/ `,
"        /___/           ",
}

artWidth := 25
leftPad := (a.width - artWidth) / 2
if leftPad < 0 {
leftPad = 0
}
pad := strings.Repeat(" ", leftPad)

hint := HelpDescStyle.Render("[?] help  [q] quit")
hintPad := a.width - lipgloss.Width(hint) - 1
if hintPad < 0 {
hintPad = 0
}

var rows []string
// top blank line with hint in the right corner
rows = append(rows, strings.Repeat(" ", hintPad)+hint)
for _, line := range lines {
rows = append(rows, pad+TitleStyle.Render(line))
}
// bottom blank line
rows = append(rows, "")
return strings.Join(rows, "\n")
}

func (a App) viewHelpBar() string {
	var parts []string
	sep := HelpSepStyle.Render("  •  ")

	addBinding := func(key, desc string) {
		parts = append(parts, HelpKeyStyle.Render(key)+" "+HelpDescStyle.Render(desc))
	}

	if a.mode == modeSearch {
		addBinding("enter/esc", "done searching")
	} else {
		addBinding("tab", "switch focus")
		addBinding("a", "add")
		addBinding("e", "edit")
		addBinding("d", "delete")
		addBinding("/", "search")
	}

	bar := strings.Join(parts, sep)
	if a.mode == modeSearch {
		bar = SearchPromptStyle.Render("/") + " " + a.search.View() + "  " + bar
	}

	return StatusBarStyle.Width(a.width).Render(bar)
}

func (a App) viewForm() string {
	return a.form.View()
}

func (a App) viewConfirmOverlay() string {
	return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center,
		a.confirm.View(a.width, a.height),
		lipgloss.WithWhitespaceForeground(colorMuted),
	)
}

func (a App) viewHelp() string {
	sections := []struct{ key, desc string }{
		{"tab / shift+tab", "Switch focus between sidebar and list"},
		{"j / k / ↑ / ↓", "Navigate items"},
		{"a", "Add new entry"},
		{"e", "Edit selected entry"},
		{"d", "Delete selected entry"},
		{"/", "Search / filter entries"},
		{"esc", "Clear search"},
		{"?", "Toggle this help"},
		{"q / ctrl+c", "Quit"},
	}

	var rows []string
	rows = append(rows, FormTitleStyle.Render(fmt.Sprintf("klyde %s — help", a.version)))
	rows = append(rows, "")

	for _, s := range sections {
		row := lipgloss.JoinHorizontal(lipgloss.Top,
			HelpKeyStyle.Width(24).Render(s.key),
			HelpDescStyle.Render(s.desc),
		)
		rows = append(rows, row)
	}

	rows = append(rows, "")
	rows = append(rows, HelpDescStyle.Render("Press ? or esc to close"))
	rows = append(rows, "")
	rows = append(rows, HelpSepStyle.Render("─────────────────────────────────────"))
	rows = append(rows, HelpDescStyle.Render("Data  ")+HelpKeyStyle.Render(a.store.Path()))
	rows = append(rows, HelpDescStyle.Render("Morten Johansen ")+HelpSepStyle.Render("| ")+HelpKeyStyle.Render("johansen.foo"))
	rows = append(rows, HelpDescStyle.Render("github.com/mojoaar/klyde"))

	content := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(2, 4).
		Render(lipgloss.JoinVertical(lipgloss.Left, rows...))

	return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, content)
}
