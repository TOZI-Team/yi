/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"Yi/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	cmd.Execute()
	//project.InitGuide(t.DefaultConfig)
}
