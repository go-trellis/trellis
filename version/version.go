package version

import (
	"fmt"

	"trellis.tech/trellis/common.v1/builder"
)

// ShowAllInfo 展示详细信息
func ShowAllInfo() {
	builder.Show()
}

// Version 版本信息
func Version() string {
	return fmt.Sprintf("%s, version: %s (branch: %s, revision: %s)",
		builder.ProgramName, builder.ProgramVersion,
		builder.ProgramBranch, builder.ProgramRevision,
	)
}

// BuildInfo returns goVersion, Author and buildTime information.
func BuildInfo() string {
	return fmt.Sprintf("(go=%s, user=%s, date=%s)",
		builder.CompilerVersion, builder.Author, builder.BuildTime)
}
