package subcmds

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vaaandark/kaleido-runtime/pkg/kmsglog"
	"github.com/vaaandark/kaleido-runtime/pkg/runm/migration"
)

func NewCreateCommand() *cobra.Command {
	// subcommand flags
	var bundle string
	var consoleSocket string
	var pidFile string
	var noPivot bool
	var noNewKeyring bool
	var preserveFds int

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create task from checkpoint",
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
			kmsglog.InfoF("subcommand flags: bundle=%s, console-socket=%s, pid-file=%s, no-pivot=%t, no-new-keyring=%t, preserve-fds=%d",
				bundle, consoleSocket, pidFile, noPivot, noNewKeyring, preserveFds)

			containerInfo := migration.NewContainerInfo(args[len(args)-1], bundle, root)
			migration, shouldMigrate, err := migration.NewMigration(containerInfo, migration.CreateAction)
			if err != nil {
				return err
			}
			if !shouldMigrate {
				return fmt.Errorf("container %s need not to be migrated", containerInfo.Id)
			}

			kmsglog.InfoF("containerInfo: %+v", containerInfo)
			kmsglog.InfoF("migration: %+v", migration)
			kmsglog.InfoF("should migrate: %t", shouldMigrate)

			if err := migration.Restore(); err != nil {
				kmsglog.InfoF("Failed to restore container %s: %v", containerInfo.Id, err)
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&bundle, "bundle", "", "path to the root of the bundle directory, defaults to the current directory")
	cmd.Flags().StringVar(&consoleSocket, "console-socket", "", "path to an AF_UNIX socket which will receive a file descriptor referencing the master end of the console's pseudoterminal")
	cmd.Flags().StringVar(&pidFile, "pid-file", "", "specify the file to write the process id to")
	cmd.Flags().BoolVar(&noPivot, "no-pivot", false, "do not use pivot root to jail process inside rootfs.  This should be used whenever the rootfs is on top of a ramdisk")
	cmd.Flags().BoolVar(&noNewKeyring, "no-new-keyring", false, "do not create a new session keyring for the container.  This will cause the container to inherit the calling processes session key")
	cmd.Flags().IntVar(&preserveFds, "preserve-fds", 0, "Pass N additional file descriptors to the container (stdio + $LISTEN_FDS + N in total)")

	return cmd
}
