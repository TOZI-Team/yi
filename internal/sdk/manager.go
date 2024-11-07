package sdk

import "Yi/pkg/types"

type Manager struct {
	sdks []types.SDKInfo
}

func NewSDKManager(sdks []types.SDKInfo) *Manager {
	return &Manager{sdks: sdks}
}

func (m Manager) GetSDKs() *[]types.SDKInfo {
	return &m.sdks
}

var GlobalSDKManger *Manager

func init() {
	GlobalSDKManger = NewSDKManager([]types.SDKInfo{{Path: "/opt/cangjie", Ver: "0.53.13", Note: "Beta Archive"}})
}
