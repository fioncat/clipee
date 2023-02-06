package daemon

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fioncat/clipee/pkg/daemon"
	"github.com/spf13/cobra"
)

var (
	flagTargets  []string
	flagNoDaemon bool
)

func init() {
	Start.PersistentFlags().StringSliceVarP(&flagTargets, "target", "t", nil, "publish targets")
	Start.PersistentFlags().BoolVarP(&flagNoDaemon, "no-daemon", "n", false, "no daemon")
}

var Start = &cobra.Command{
	Use:   "start <listen-addr> [-t target]... [-n]",
	Short: "Start share",

	RunE: func(_ *cobra.Command, args []string) error {
		if !flagNoDaemon {
			d, err := daemon.New()
			if err != nil {
				return fmt.Errorf("failed to init daemon: %v", err)
			}
			return d.Start(func() error {
				return start(args[0])
			})
		}

		return start(args[0])
	},

	Args: cobra.ExactArgs(1),
}

var Stop = &cobra.Command{
	Use:   "stop",
	Short: "Stop daemon",

	RunE: func(_ *cobra.Command, _ []string) error {
		d, err := daemon.New()
		if err != nil {
			return err
		}
		return d.Stop()
	},
}

var Status = &cobra.Command{
	Use:   "status",
	Short: "Show daemon status",

	RunE: func(_ *cobra.Command, _ []string) error {
		d, err := daemon.New()
		if err != nil {
			return err
		}
		return d.ShowStatus()
	},
}

var Logs = &cobra.Command{
	Use:   "logs",
	Short: "Show daemon logs",

	DisableFlagParsing: true,

	RunE: func(_ *cobra.Command, args []string) error {
		d, err := daemon.New()
		if err != nil {
			return err
		}
		path := d.LogPath()
		args = append(args, path)
		cmd := exec.Command("tail", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		return cmd.Run()
	},
}

func start(addr string) error {
	if len(flagTargets) > 0 {
		go Client(flagTargets)
	}
	return Server(addr)
}
