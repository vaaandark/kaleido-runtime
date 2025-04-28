package main

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/vaaandark/kaleido-runtime/cmd/runm/subcmds"
	"github.com/vaaandark/kaleido-runtime/pkg/kmsglog"
)

const (
	runc = "runc"
)

func main() {
	runcArgs := os.Args[1:]

	rootCmd := &cobra.Command{
		Use:              "runm",
		Short:            "runm is a container runtime wrapper",
		TraverseChildren: true,
		SilenceErrors:    true,
		SilenceUsage:     true,
	}

	// Define root flags
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug logging")
	rootCmd.PersistentFlags().String("log", "", "set the log file to write runc logs to (default is '/dev/stderr')")
	rootCmd.PersistentFlags().String("log-format", "text", "set the log format ('text' (default), or 'json')")
	rootCmd.PersistentFlags().String("root", "", "root directory for storage of container state (this should be located in tmpfs)")
	rootCmd.PersistentFlags().String("criu", "criu", "path to the criu binary used for checkpoint and restore")
	rootCmd.PersistentFlags().Bool("systemd-cgroup", false, "enable systemd cgroup support")
	rootCmd.PersistentFlags().String("rootless", "auto", "ignore cgroup permission errors ('true', 'false', or 'auto')")

	rootCmd.AddCommand(subcmds.NewCreateCommand())
	rootCmd.AddCommand(subcmds.NewKillCommand())

	if cmd, _, err := rootCmd.Find(runcArgs); err == nil &&
		(cmd.Use == "create" || cmd.Use == "kill") {
		if err := rootCmd.Execute(); err != nil {
			kmsglog.InfoF("Failed to execute runm: %v", err)
		} else {
			return
		}
	}

	runcCmd := exec.Command(runc, runcArgs...)
	runcCmd.Stdout = os.Stdout
	runcCmd.Stderr = os.Stderr
	kmsglog.InfoF("Running runc %v", append([]string{runc}, runcArgs...))
	if err := runcCmd.Run(); err != nil {
		kmsglog.WarnF("Failed to run runc. If it is a checkpoint restoring, just ignore it")
	}
}
