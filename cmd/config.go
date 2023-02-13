package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var Config = &cobra.Command{
	Use:   "config",
	Short: "Edit your config file",

	RunE: func(_ *cobra.Command, _ []string) error {
		path := os.Getenv("CLIPEE_CONFIG")
		if path == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			path = filepath.Join(homeDir, ".config", "clipee", "config.yaml")
		}

		editor, err := selectEditor()
		if err != nil {
			return err
		}

		cmd := exec.Command(editor, path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		return cmd.Run()
	},
}

var editors = []string{
	"nvim", "vim", "vi",
}

func selectEditor() (string, error) {
	if e := os.Getenv("CLIPEE_EDITOR"); e != "" {
		return e, nil
	}
	for _, editor := range editors {
		if commandExists(editor) {
			return editor, nil
		}
	}
	return "", errors.New("no editor to use")
}

func commandExists(name string) bool {
	args := []string{
		"-c",
		fmt.Sprintf("command -v %s", name),
	}
	cmd := exec.Command("bash", args...)
	return cmd.Run() == nil
}
