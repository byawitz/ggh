package theme

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type Print int

const (
	PrintConfig = iota
	PrintHistory
)

func PrintTable(rows []table.Row, p Print) string {
	columns := []table.Column{
		{Title: "Name", Width: 10},
		{Title: "Host", Width: 15},
		{Title: "Port", Width: 10},
		{Title: "User", Width: 10},
		{Title: "Key", Width: 10},
	}

	if p == PrintHistory {
		columns = append(columns, table.Column{Title: "Last login", Width: 15})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithStyles(table.Styles{
			Header:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")),
			Selected: lipgloss.NewStyle(),
		}),
		table.WithHeight(len(rows)+1),
	)

	return BaseStyle.Render(t.View())
}
