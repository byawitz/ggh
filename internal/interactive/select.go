package interactive

import (
	"fmt"
	"github.com/byawitz/ggh/internal/config"
	"github.com/byawitz/ggh/internal/theme"
	"math"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Selecting int

const (
	SelectConfig Selecting = iota
	SelectHistory
)

type model struct {
	table  table.Model
	choice config.SSHConfig
	what   Selecting
	exit   bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
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
	if what == SelectConfig {
		return config.SSHConfig{
			Host: row[1],
			Port: row[2],
			User: row[3],
			Key:  row[4],
		}
	}

	return config.SSHConfig{
		Host: row[0],
		Port: row[1],
		User: row[2],
		Key:  row[3],
	}
}

func (m model) View() string {
	if m.choice.Host != "" || m.exit {
		return ""
	}
	return theme.BaseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n"
}

func Select(rows []table.Row, what Selecting) config.SSHConfig {
	var columns []table.Column
	if what == SelectConfig {
		columns = append(columns, []table.Column{
			{Title: "Name", Width: 15},
			{Title: "Host", Width: 15},
			{Title: "Port", Width: 10},
			{Title: "User", Width: 10},
			{Title: "Key", Width: 10},
		}...)
	}

	if what == SelectHistory {
		columns = append(columns, []table.Column{
			{Title: "Host", Width: 15},
			{Title: "Port", Width: 10},
			{Title: "User", Width: 10},
			{Title: "Key", Width: 10},
			{Title: "Last login", Width: 15},
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
