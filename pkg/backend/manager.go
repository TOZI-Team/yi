package backend

import (
	cjpackage "yi/pkg/package"
)

type Manager struct {
	bs []cjpackage.ProjectBackendInterface
}

func NewManager() *Manager {
	return new(Manager)
}

func (m *Manager) AddBackend(backend cjpackage.ProjectBackendInterface) {
	m.bs = append(m.bs, backend)
}

func (m *Manager) GetBackend() cjpackage.ProjectBackendInterface {
	return m.bs[0]
}

var GlobalBackendManager *Manager

func init() {
	GlobalBackendManager = NewManager()
	cjpackage.SDKBackendManager = GlobalBackendManager
}
