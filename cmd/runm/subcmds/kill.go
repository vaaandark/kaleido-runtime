package subcmds

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaaandark/kaleido-runtime/pkg/kmsglog"
	"github.com/vaaandark/kaleido-runtime/pkg/runm/migration"
	"github.com/vaaandark/kaleido-runtime/pkg/runtimeutils"
)

func IsDirExist(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return fileInfo.IsDir()
}

func NewKillCommand() *cobra.Command {
	// subcommand flags
	var all bool

	cmd := &cobra.Command{
		Use:   "kill",
		Short: "Convert kill signal to checkpointing tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			debug, _ := cmd.Root().Flags().GetBool("debug")
			log, _ := cmd.Root().Flags().GetString("log")
			logFormat, _ := cmd.Root().Flags().GetString("log-format")
			root, _ := cmd.Root().Flags().GetString("root")
			criu, _ := cmd.Root().Flags().GetString("criu")
			systemdCgroup, _ := cmd.Root().Flags().GetBool("systemd-cgroup")
			rootless, _ := cmd.Root().Flags().GetString("rootless")
			kmsglog.InfoF("root command flags: debug=%t, log=%s, log-format=%s, root=%s, criu=%s, systemd-cgroup=%t, rootless=%s",
				debug, log, logFormat, root, criu, systemdCgroup, rootless)

			containerId := args[len(args)-2]
			signal := args[len(args)-1]
			state, err := runtimeutils.LoadLibcontainerState(root, containerId)
			if err != nil {
				kmsglog.ErrorF("Fail to load libcontainer state: %v", err)
				return err
			}
			bundle := runtimeutils.GetBundleFromState(state)

			containerInfo := migration.NewContainerInfo(containerId, bundle, root)
			migration, shouldMigrate, err := migration.NewMigration(containerInfo, migration.KillAction)
			if err != nil {
				return err
			}

			if !IsDirExist(migration.CriuRoot()) {
				return fmt.Errorf("criu root directory %s not exist", migration.CriuRoot())
			}

			if !shouldMigrate {
				return fmt.Errorf("container %s need not to be migrated", containerInfo.Id)
			}

			kmsglog.InfoF("containerInfo: %+v", containerInfo)
			kmsglog.InfoF("migration: %+v", migration)
			kmsglog.InfoF("should migrate: %t", shouldMigrate)

			// ignore signal kill
			if signal != "9" {
				return nil
			}

			if err := migration.Checkpoint(); err != nil {
				kmsglog.InfoF("Failed to checkpoint container %s: %v", containerInfo.Id, err)
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, " send the specified signal to all processes inside the container")

	return cmd
}
