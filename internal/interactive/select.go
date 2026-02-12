package interactive

import (
	"fmt"
	"github.com/byawitz/ggh/internal/config"
	"github.com/byawitz/ggh/internal/history"
	"github.com/byawitz/ggh/internal/settings"
	"github.com/byawitz/ggh/internal/theme"
	"math"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Selecting int

const (
	SelectConfig Selecting = iota
	SelectHistory
)

const (
	MarginWidth            = 3
	MarginHeight           = 5
	MinimumTableWidth      = 3
	ContentExtraMargin     = 12
	PreferredKeyExtraWidth = 15
	MaxKeyExtraWidth       = 30
)

type model struct {
	table        table.Model
	choice       config.SSHConfig
	what         Selecting
	exit         bool
	windowWidth  int
	windowHeight int
	settings     settings.Settings
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	// 1. Handle window resize events
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height - MarginHeight

		widthForTable := max(m.windowWidth-MarginWidth, MinimumTableWidth)
		// Extra margin for content
		widthForTableContent := widthForTable - ContentExtraMargin

		cols := m.table.Columns()

		switch m.what {
		// SELECT CONFIG
		case SelectConfig:
			// columns = [Name, Host, Port, User, Key]
			// base widths = 15,20,5,10,10 = total 60
			baseWidths := []int{15, 20, 5, 10, 10}
			const totalBase = 60

			if widthForTableContent >= totalBase {
				leftover := widthForTableContent - totalBase
				leftoverForKey := 0
				leftoverForName := 0

				for leftover > 0 {
					if leftoverForKey < PreferredKeyExtraWidth {
						leftoverForKey++
						leftover--
					} else if leftoverForKey < MaxKeyExtraWidth && leftover > 1 {
						leftoverForName++
						leftoverForKey++
						leftover -= 2
					} else {
						leftoverForName++
						leftover--
					}
				}

				cols[0].Width = baseWidths[0] + leftoverForName // Name
				cols[1].Width = baseWidths[1]                   // Host
				cols[2].Width = baseWidths[2]                   // Port
				cols[3].Width = baseWidths[3]                   // User
				cols[4].Width = baseWidths[4] + leftoverForKey  // Key
			} else {
				// Scale all columns proportionally
				ratio := float64(widthForTableContent) / float64(totalBase)
				for i := range cols {
					w := max(int(math.Round(float64(baseWidths[i])*ratio)), 1)
					cols[i].Width = w
				}
			}

		// SELECT HISTORY
		case SelectHistory:
			// columns = [Name,Host,Port,User,Key,Last login]
			// base widths = 10,20,5,10,0,15 = total 60
			baseWidths := []int{10, 20, 5, 10, 0, 15}
			const totalBase = 60

			if widthForTableContent >= totalBase {
				leftover := widthForTableContent - totalBase
				leftoverForKey := 0
				leftoverForName := 0

				for leftover > 0 {
					if leftoverForKey < PreferredKeyExtraWidth {
						leftoverForKey++
						leftover--
					} else if leftoverForKey < MaxKeyExtraWidth && leftover > 1 {
						leftoverForName++
						leftoverForKey++
						leftover -= 2
					} else {
						leftoverForName++
						leftover--
					}
				}

				cols[0].Width = baseWidths[0] + leftoverForName // Name
				cols[1].Width = baseWidths[1]                   // Host
				cols[2].Width = baseWidths[2]                   // Port
				cols[3].Width = baseWidths[3]                   // User
				cols[4].Width = baseWidths[4] + leftoverForKey  // Key
				cols[5].Width = baseWidths[5]                   // Last login
			} else {
				// Not enough space → scale all columns proportionally
				ratio := float64(widthForTableContent) / float64(totalBase)
				for i := range cols {
					w := max(int(math.Round(float64(baseWidths[i])*ratio)), 1)
					cols[i].Width = w
				}
			}
		}

		// Apply the new widths
		m.table.SetColumns(cols)
		m.table.SetWidth(widthForTable)
		m.settings = settings.FetchWithDefaultFile()
		if m.settings.Fullscreen {
			// if fullscreen, let the table be as tall as the terminal
			m.table.SetHeight(m.windowHeight)
			return m, tea.EnterAltScreen
		} else {
			// if not fullscreen, set the height to a minimum of 8 rows
			m.table.SetHeight(int(math.Min(8, float64(len(m.table.Rows())+2))))
			return m, tea.ExitAltScreen
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "d":
			history.RemoveByIP(m.table.SelectedRow())

			rows := slices.Delete(m.table.Rows(), m.table.Cursor(), m.table.Cursor()+1)
			m.table.SetRows(rows)

			m.table, cmd = m.table.Update("") // Overrides default `d` behavior
			return m, cmd
		case "w":
			// toggle fullscreen mode
			newsettings := m.settings
			newsettings.Fullscreen = !m.settings.Fullscreen
			if s, err := settings.Save(newsettings); err == nil && s != nil {
				m.settings = *s
				if m.settings.Fullscreen {
					// if fullscreen, let the table be as tall as the terminal
					m.table.SetHeight(m.windowHeight)
					return m, tea.EnterAltScreen
				} else {
					// if not fullscreen, set the height to a minimum of 8 rows
					m.table.SetHeight(int(math.Min(8, float64(len(m.table.Rows())+2))))
					return m, tea.ExitAltScreen
				}
			}

			// If we can't save the settings, do nothing
			return m, nil
		case "q", "ctrl+c", "esc":
			m.exit = true
			return m, tea.Quit
		case "enter":
			m.choice = setConfig(m.table.SelectedRow(), m.what)
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func setConfig(row table.Row, what Selecting) config.SSHConfig {
	return config.SSHConfig{
		Host: row[1],
		Port: row[2],
		User: row[3],
		Key:  row[4],
	}
}

func (m model) View() string {
	if m.choice.Host != "" || m.exit {
		return ""
	}
	return theme.BaseStyle.Render(m.table.View()) + "\n  " + m.HelpView() + "\n"
}

func Select(rows []table.Row, what Selecting) config.SSHConfig {
	var columns []table.Column
	if what == SelectConfig {
		columns = append(columns, []table.Column{
			{Title: "Name"},
			{Title: "Host"},
			{Title: "Port"},
			{Title: "User"},
			{Title: "Key"},
		}...)
	}

	if what == SelectHistory {
		columns = append(columns, []table.Column{
			{Title: "Name"},
			{Title: "Host"},
			{Title: "Port"},
			{Title: "User"},
			{Title: "Key"},
			{Title: "Last login"},
		}...)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(int(math.Min(8, float64(len(rows)+1)))),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(false)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)

	t.SetStyles(s)

	p := tea.NewProgram(model{table: t, what: what})
	m, err := p.Run()
	if err != nil {
		fmt.Println("error while running the interactive selector, ", err)
		os.Exit(1)
	}
	// Assert the final tea.Model to our local model and print the choice.
	if m, ok := m.(model); ok {
		if m.choice.Host != "" {
			return m.choice
		}
		if m.exit {
			os.Exit(0)
		}
	}

	return config.SSHConfig{}
}
func (m model) HelpView() string {

	km := table.DefaultKeyMap()

	var b strings.Builder

	b.WriteString(generateHelpBlock(km.LineUp.Help().Key, km.LineUp.Help().Desc, true))
	b.WriteString(generateHelpBlock(km.LineDown.Help().Key, km.LineDown.Help().Desc, true))

	if m.what == SelectHistory {
		b.WriteString(generateHelpBlock("d", "delete", true))
	}

	b.WriteString(generateHelpBlock("w", "full/windowed", true))
	b.WriteString(generateHelpBlock("q/esc", "quit", false))

	return b.String()
}

func generateHelpBlock(key, desc string, withSep bool) string {
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#909090",
		Dark:  "#626262",
	})

	descStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#B2B2B2",
		Dark:  "#4A4A4A",
	})

	sepStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#DDDADA",
		Dark:  "#3C3C3C",
	})

	sep := sepStyle.Inline(true).Render(" • ")

	str := keyStyle.Inline(true).Render(key) +
		" " +
		descStyle.Inline(true).Render(desc)

	if withSep {
		str += sep
	}

	return str
}
