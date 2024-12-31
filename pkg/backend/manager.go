package backend

import t "yi/pkg/types"

type Manager struct {
	bs []t.ProjectBackendInterface
}

func NewManager() *Manager {
	return new(Manager)
}

func (m *Manager) AddBackend(backend t.ProjectBackendInterface) {
	m.bs = append(m.bs, backend)
}

func (m *Manager) GetBackend() t.ProjectBackendInterface {
	return m.bs[0]
}

var GlobalBackendManager *Manager

func init() {
	GlobalBackendManager = NewManager()
}
