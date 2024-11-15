/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	log "github.com/sirupsen/logrus"
	"yi/cmd"
)

func main() {
	log.SetLevel(log.DebugLevel)
	cmd.Execute()
	//project.InitGuide(t.DefaultConfig)
}
