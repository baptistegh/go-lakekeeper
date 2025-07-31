package main

import (
	"os"

	"github.com/baptistegh/go-lakekeeper/cmd/lkctl/commands"
)

func main() {
	command := commands.NewCommand()

	command.SilenceErrors = true
	command.SilenceUsage = true

	err := command.Execute()

	if err != nil {
		os.Exit(1)
	}
}
