/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"Yi/cmd/project"
	sdkCmd "Yi/cmd/sdk"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "Yi",
	Short:   "Cangjie package manager",
	Version: "0.1.0-dev-snapshot",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		//
		//
		//
		//          ,---,          ,--,
		//         /_ ./|        ,--.'|
		//   ,---, |  ' :        |  |,
		//  /___/ \.  : |        `--'_
		//
		//
		//
		//
		//
		//
		//
		//          `--`          ---`-'
		//
		println()
		println("          ,---,          ,--,    ")
		println("         /_ ./|        ,--.'|    ")
		println("   ,---, |  ' :        |  |,     ")
		println("   .  \\  \\ ,' '        ,' ,'|    ")
		println("    \\  ;  `  ,'        '  | |    ")
		println("     \\  \\    '         |  | :    ")
		println("      '  \\   |         '  : |__  ")
		println("       \\  ;  ;         |  | '.'| ")
		println("        :  \\  \\        ;  :    ; ")
		println("         \\  ' ;        |  ,   /  ")
		println("          `--`          ---`-'   ")
		println("=================================")
		//println(cmd.VersionTemplate())
		fmt.Printf("Yi %s\n", cmd.Version)
		err := cmd.Help()
		if err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.AddCommand(project.NewCommand)
	rootCmd.AddCommand(project.BuildCommand)
	rootCmd.AddCommand(sdkCmd.Command)
	rootCmd.AddCommand(project.InitCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
