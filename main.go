package main

import (
	"github.com/hichtakk/nsxctl/cmd"
)

func main() {
	rootCmd := cmd.GetCmdRoot()
	rootCmd.Execute()
	return
}
