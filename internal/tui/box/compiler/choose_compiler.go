package compiler

import (
	"Yi/internal/tui/modules"
	t "Yi/pkg/types"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/sirupsen/logrus"
	"os"
)

func ChooseCompiler(sdks *[]t.SDKInfo) t.SDKInfo {
	var rows []table.Row
	for _, sdk := range *sdks {
		rows = append(rows, table.Row{sdk.Ver, sdk.Path, sdk.Note})
	}
	columns := []table.Column{
		{Title: "Version", Width: 10},
		{Title: "Path", Width: 25},
		{Title: "Note", Width: 15},
	}

	model := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	model.SetStyles(s)

	m := modules.NewChooseModel(model)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	for _, sdk := range *sdks {
		if sdk.Path == m.ChooseIndex[1] {
			//log.Info("Found sdk")
			return sdk
		}
	}
	log.Warnf("Can not find this SDK. Use default: %s", (*sdks)[0].Ver)
	return (*sdks)[0]
}
