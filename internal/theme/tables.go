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

func CalculateColumnWidths(rows []table.Row, columns []table.Column) []table.Column {
	padding := 2
	maxWidths := make([]int, len(columns))
	for i, col := range columns {
		maxWidths[i] = len(col.Title) + padding
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(maxWidths) {
				width := len(cell) + padding
				if width > maxWidths[i] {
					maxWidths[i] = width
				}
			}
		}
	}

	for i := range columns {
		columns[i].Width = maxWidths[i]
	}

	return columns
}

func PrintTable(rows []table.Row, p Print) string {
	columns := []table.Column{
		{Title: "Name"},
		{Title: "Host"},
		{Title: "Port"},
		{Title: "User"},
		{Title: "Key"},
	}

	if p == PrintHistory {
		columns = append(columns, table.Column{Title: "Last login"})
	}

	columns = CalculateColumnWidths(rows, columns)

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
