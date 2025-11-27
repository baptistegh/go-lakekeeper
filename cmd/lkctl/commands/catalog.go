package commands

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCatalogCmd(_ *clientOptions) *cobra.Command {
	command := cobra.Command{
		Use:   "catalog",
		Short: "Interacts with catalogs (not implemented)",
		Run: func(_ *cobra.Command, _ []string) {
			logrus.Fatal("catalog command is not implemented")
		},
	}

	return &command
}
