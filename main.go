package main

import (
	"fmt"

	"github.com/fioncat/clipee/cmd/daemon"
	"github.com/fioncat/clipee/config"
	"github.com/fioncat/clipee/pkg/clipboard"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "clipee",

	SilenceErrors: true,
	SilenceUsage:  true,

	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		err := config.Init()
		if err != nil {
			return err
		}

		err = clipboard.Init()
		if err != nil {
			return fmt.Errorf("failed to init clipboard driver: %v", err)
		}
		return nil
	},

	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

func main() {
	Cmd.AddCommand(daemon.Start, daemon.Stop, daemon.Status, daemon.Logs)

	err := Cmd.Execute()
	if err != nil {
		logrus.Fatal(err)
	}
}
