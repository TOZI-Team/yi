/*
Copyright © 2024 Yan Mingshu <yms_hi@Outlook.com>
*/
package main

import (
	log "github.com/sirupsen/logrus"
	"yi/cmd"
	"yi/pkg/backend"
	cjpmPackage "yi/pkg/backend/cjpm/package"
)

func main() {
	backend.GlobalBackendManager.AddBackend(cjpmPackage.CJPMProjectBackend{})
	log.SetLevel(log.DebugLevel)
	cmd.Execute()
	//project.InitGuide(t.DefaultConfig)
}
