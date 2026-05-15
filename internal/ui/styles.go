package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Palette
	colorPrimary  = lipgloss.Color("#7C3AED") // violet
	colorAccent   = lipgloss.Color("#A78BFA") // light violet
	colorMuted    = lipgloss.Color("#6B7280") // gray
	colorText     = lipgloss.Color("#F9FAFB") // near-white
	colorSubtle   = lipgloss.Color("#D1D5DB") // light gray
	colorBorder   = lipgloss.Color("#374151") // dark gray border
	colorSelected = lipgloss.Color("#1F1035") // dark violet bg
	colorDanger   = lipgloss.Color("#EF4444") // red for delete

	// App chrome
	TitleStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Padding(0, 1)

	// Sidebar
	SidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), false, true, false, false).
			BorderForeground(colorBorder).
			Padding(1, 1).
			Width(20)

	SidebarItemStyle = lipgloss.NewStyle().
				Foreground(colorSubtle).
				Padding(0, 1)

	SidebarActiveStyle = lipgloss.NewStyle().
				Foreground(colorAccent).
				Bold(true).
				Padding(0, 1)

	SidebarFocusBorderStyle = lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder(), false, true, false, false).
					BorderForeground(colorPrimary).
					Padding(1, 1).
					Width(20)

	// List
	ListStyle = lipgloss.NewStyle().
			Padding(1, 2)

	ListItemStyle = lipgloss.NewStyle().
			Foreground(colorText)

	ListSelectedStyle = lipgloss.NewStyle().
				Background(colorSelected).
				Foreground(colorAccent).
				Bold(true)

	// Column headers
	HeaderStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Bold(true)

	// Form
	FormLabelStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	FormActiveInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorPrimary).
				Padding(0, 1)

	FormInputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1)

	FormTitleStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true).
			MarginBottom(1)

	// Confirm dialog
	ConfirmBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorDanger).
			Padding(1, 3).
			Align(lipgloss.Center)

	ConfirmTitleStyle = lipgloss.NewStyle().
				Foreground(colorDanger).
				Bold(true).
				MarginBottom(1)

	ConfirmYesStyle = lipgloss.NewStyle().
			Background(colorDanger).
			Foreground(colorText).
			Bold(true).
			Padding(0, 2)

	ConfirmNoStyle = lipgloss.NewStyle().
			Foreground(colorSubtle).
			Padding(0, 2)

	// Help bar
	HelpKeyStyle  = lipgloss.NewStyle().Foreground(colorAccent)
	HelpDescStyle = lipgloss.NewStyle().Foreground(colorMuted)
	HelpSepStyle  = lipgloss.NewStyle().Foreground(colorBorder)

	// Search
	SearchPromptStyle = lipgloss.NewStyle().Foreground(colorAccent)
	SearchStyle       = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorPrimary).
				Padding(0, 1)
)
