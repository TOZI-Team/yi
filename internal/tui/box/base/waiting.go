package tBox

import (
	tea "github.com/charmbracelet/bubbletea"
	"yi/internal/tui/modules"
	t "yi/pkg/types"
)

type BaseWaitingBox struct {
	m *modules.WaitingModel
	f func(chan t.WaitingMessage)
}

func (box BaseWaitingBox) Run() error {
	c := make(chan t.WaitingMessage, 1)
	box.m = modules.NewWaitingModel("Working", c)
	go box.f(c)

	if _, err := tea.NewProgram(box.m).Run(); err != nil {
		return err
	}
	return nil
}

func NewWaitingBox(f func(chan t.WaitingMessage)) *BaseWaitingBox {
	return &BaseWaitingBox{
		f: f,
		m: nil,
	}
}

type JobsModel struct {
	jobs []t.Job
}

func (m JobsModel) Run(c chan t.WaitingMessage) error {
	for _, job := range m.jobs {
		c <- t.WaitingMessage{
			Message: job.Name,
			Statue:  t.Waiting,
		}
		err := job.Run()
		if err != nil {
			c <- t.WaitingMessage{
				Message: err.Error(),
				Statue:  t.Err,
			}
		}
		return err
	}
	c <- t.WaitingMessage{
		Statue: t.Success,
	}
	return nil
}
