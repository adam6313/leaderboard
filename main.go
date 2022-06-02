package main

import (
	"leaderboard/cmd"
	"leaderboard/config"
)

var (
	// VERSION -
	VERSION string
)

func main() {
	config.C.Version = VERSION

	cmd.Execute()
}
