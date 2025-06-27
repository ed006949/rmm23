package l

var (
	gitCommit string
)
var (
	control = &ControlType{
		Name:      "",
		Config:    "",
		DryRun:    DoDryRun,
		Mode:      Init,
		Verbosity: Informational,
		GitCommit: gitCommitValue(gitCommit),
	}
)
var (
	dryRunDescription = map[dryRunFlag]string{
		NoDryRun: "false",
		DoDryRun: "true",
	}
)
var (
	modeDescription = map[modeValue]string{
		Init:   "init",
		Deploy: "deploy",
		CLI:    "cli",
		Daemon: "daemon",
	}
)
