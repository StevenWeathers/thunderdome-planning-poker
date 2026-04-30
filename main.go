package main

import (
	_ "time/tzdata"

	"github.com/StevenWeathers/thunderdome-planning-poker/cmd"
)

var (
	version = "dev"
)

func main() {
	cmd.Execute(version)
}
