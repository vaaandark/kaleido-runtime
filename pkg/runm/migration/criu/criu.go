package criu

type Options struct {
	LazyPages         bool
	ImagePath         string
	WorkPath          string
	TcpEstablished    bool
	AutoDedup         bool
	ManageCgroupsMode string
	NoSubreaper       bool
	Detach            bool
}

func (o *Options) BuildRestoreArgs() []string {
	restoreArgs := []string{}
	appendBoolArg := func(flagName string, value bool) {
		if value {
			restoreArgs = append(restoreArgs, flagName)
		}
	}
	appendStringArg := func(flagName, value string) {
		if value != "" {
			restoreArgs = append(restoreArgs, flagName, value)
		}
	}
	appendBoolArg("--lazy-pages", o.LazyPages)
	appendStringArg("--image-path", o.ImagePath)
	appendStringArg("--work-path", o.WorkPath)
	appendBoolArg("--tcp-established", o.TcpEstablished)
	appendBoolArg("--auto-dedup", o.AutoDedup)
	appendStringArg("--manage-cgroups-mode", o.ManageCgroupsMode)
	appendBoolArg("--no-subreaper", o.NoSubreaper)
	appendBoolArg("--detach", o.Detach)
	return restoreArgs
}
