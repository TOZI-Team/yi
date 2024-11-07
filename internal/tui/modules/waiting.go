package modules

import (
	t "Yi/pkg/types"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WaitingModel 提供等待任务的TUI
type WaitingModel struct {
	spinner  spinner.Model
	err      error
	statue   t.WaitingMessage
	quitting bool
	statueC  chan t.WaitingMessage
}

func (m WaitingModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m WaitingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(m.statueC) == 1 {
		tsatue := <-m.statueC
		if tsatue.Message != "" {
			m.statue.Message = tsatue.Message
		}
		m.statue.Statue = tsatue.Statue
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			return m, nil
		}
	case error:
		m.err = msg
		return m, nil
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

}

func (m WaitingModel) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	var icon string
	switch m.statue.Statue {
	case t.Waiting:
		icon = m.spinner.View()
	case t.Err:
		icon = "\ue654"
	default:
		icon = m.spinner.View()
	}
	str := fmt.Sprintf("\n%s %s", icon, m.statue.Message)
	if m.statue.Statue == t.Success {
		return str + "\n"
	}
	return str
}

func NewWaitingModel(text string, c chan t.WaitingMessage) WaitingModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return WaitingModel{
		spinner: s,
		statue: t.WaitingMessage{
			Message: text,
			Statue:  t.Waiting,
		},
	}
}
