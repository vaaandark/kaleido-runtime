package migration

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/vaaandark/kaleido/pkg/runtimeutils"
	"github.com/vaaandark/kaleido/pkg/types"
)

func ShouldMigrateByBundle(bundle string) bool {
	spec, err := runtimeutils.LoadRuntimeSpecByBundle(bundle)
	if err != nil {
		return false
	}
	return ShouldMigrate(spec)
}

func ShouldMigrate(spec *specs.Spec) bool {
	if val, exist := spec.Annotations[types.MigrationAnnotation]; exist && val == "true" {
		return true
	}
	return false
}
