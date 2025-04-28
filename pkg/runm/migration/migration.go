package migration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/vaaandark/kaleido-runtime/pkg/kmsglog"
	"github.com/vaaandark/kaleido-runtime/pkg/runm/migration/criu"
	"github.com/vaaandark/kaleido-runtime/pkg/runtimeutils"
	"github.com/vaaandark/kaleido-runtime/pkg/types"
)

const (
	runc = "runc"
)

type MigrationAction string

const (
	CreateAction = "create"
	KillAction   = "kill"
)

const (
	containerTypeAnnotation = "io.kubernetes.cri.container-type"
	containerTypeSandbox    = "sandbox"
	containerTypeContainer  = "container"
	containerNameAnnotation = "io.kubernetes.cri.container-name"
	containerIdPattern      = `[a-zA-Z0-9][a-zA-Z0-9_.-]+`
)

var (
	containerIdRegex = regexp.MustCompile(containerIdPattern)
)

type ContainerInfo struct {
	Id     string
	Bundle string
	Root   string
}

func NewContainerInfo(id, bundle, root string) *ContainerInfo {
	return &ContainerInfo{
		Id:     containerIdRegex.FindString(id),
		Bundle: bundle,
		Root:   root,
	}
}

type Migration struct {
	SourcePodUid  string
	ContainerName string
	ContainerSpec *specs.Spec
	ContainerInfo *ContainerInfo
}

const (
	criuCheckpointPath = "/var/lib/criu/checkpoints"
)

func podCheckpointPath(podUid string) string {
	return filepath.Join(criuCheckpointPath, podUid)
}

func containerCheckpointPath(podUid, containerName string) string {
	return filepath.Join(podCheckpointPath(podUid), containerName)
}

func NewMigration(containerInfo *ContainerInfo, action MigrationAction) (info *Migration, shouldMigrate bool, err error) {
	spec, err := loadRuntimeSpecFromBundle(containerInfo.Bundle)
	if err != nil {
		return nil, false, fmt.Errorf("failed to load runtime spec: %w", err)
	}

	var podUid string
	var exist bool
	switch action {
	case CreateAction:
		podUid, exist = spec.Annotations[types.MigrationUidAnnotation]
	case KillAction:
		podUid, exist = spec.Annotations[types.SourceMigrationUidAnnotation]
	}
	if !exist {
		return nil, false, nil
	}

	containerType, exist := spec.Annotations[containerTypeAnnotation]
	if !exist || containerType != containerTypeContainer {
		return nil, false, fmt.Errorf("container type not found in annotations or not container type")
	}

	containerName, exist := spec.Annotations[containerNameAnnotation]
	if !exist {
		return nil, false, fmt.Errorf("container name not found in annotations")
	}

	return &Migration{
		SourcePodUid:  podUid,
		ContainerName: containerName,
		ContainerSpec: spec,
		ContainerInfo: containerInfo,
	}, true, nil
}

func loadRuntimeSpecFromBundle(bundle string) (*specs.Spec, error) {
	path := filepath.Join(bundle, "config.json")
	return runtimeutils.LoadRuntimeSpec(path)
}

func (m Migration) CriuRoot() string {
	return containerCheckpointPath(m.SourcePodUid, m.ContainerName)
}

func (m Migration) Checkpoint() error {
	criuRoot := m.CriuRoot()
	var args []string
	for _, arg := range os.Args[1 : len(os.Args)-1] {
		if arg == "kill" {
			args = append(args, "checkpoint")
			args = append(args, "--image-path", criuRoot)
		} else {
			args = append(args, arg)
		}
	}
	kmsglog.InfoF("checkpoint runc command: runc %v", args)

	runcCmd := exec.Command(runc, args...)
	runcCmd.Stdout = os.Stdout
	runcCmd.Stderr = os.Stdin

	if err := runcCmd.Run(); err != nil {
		return fmt.Errorf("failed to run checkpoint command: %w", err)
	}

	return nil
}

func (m Migration) Restore() error {
	criuRoot := m.CriuRoot()
	options := criu.Options{
		LazyPages:         false,
		ImagePath:         criuRoot,
		WorkPath:          criuRoot,
		TcpEstablished:    true,
		AutoDedup:         true,
		ManageCgroupsMode: "full",
		NoSubreaper:       true,
		Detach:            true,
	}
	restoreArgs := options.BuildRestoreArgs()

	var args []string
	for _, arg := range os.Args[1:] {
		if arg == "create" {
			args = append(args, "restore")
			args = append(args, restoreArgs...)
		} else {
			args = append(args, arg)
		}
	}

	kmsglog.InfoF("restore runc command: runc %v", args)

	runcCmd := exec.Command(runc, args...)
	runcCmd.Stdout = os.Stdout
	runcCmd.Stderr = os.Stdin

	if err := runcCmd.Run(); err != nil {
		return fmt.Errorf("failed to run restore command: %w", err)
	}

	return nil
}
