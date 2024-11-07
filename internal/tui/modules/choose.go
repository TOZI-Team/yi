package modules

import (
	"errors"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/sirupsen/logrus"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type ChooseModel struct {
	table       table.Model
	ChooseIndex []string
	err         error
}

func (m ChooseModel) Init() tea.Cmd { return nil }

func (m ChooseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "ctrl+c":
			m.err = errors.New("quite")
		case "enter":
			copy(m.ChooseIndex, m.table.SelectedRow())
			return m, tea.Quit
			//return m, tea.Batch(
			//	tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			//)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m ChooseModel) View() string {
	if m.err != nil {
		log.Exit(1)
	}
	//TODO 支持选择后隐藏
	//m.ChooseIndex = m.table.SelectedRow()
	return baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n"
}

func NewChooseModel(m table.Model) ChooseModel {
	return ChooseModel{table: m, ChooseIndex: make([]string, len(m.Columns()))}
}
