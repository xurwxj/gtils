package gop

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	// 初始化为 unknown，如果编译时没有传入这些值，则为 unknown
	GitCommitLog   = "unknown"
	GitStatus      = "unknown"
	BuildTime      = "unknown"
	BuildGoVersion = "unknown"
	Version        = "unknown"
)

// 返回单行格式
func StringifySingleLine() string {
	return fmt.Sprintf("Version=%s. GitCommitLog=%s. GitStatus=%s. BuildTime=%s. GoVersion=%s. runtime=%s/%s.",
		Version, GitCommitLog, GitStatus, BuildTime, BuildGoVersion, runtime.GOOS, runtime.GOARCH)
}

// 返回多行格式
func StringifyMultiLine() string {
	return fmt.Sprintf("Version=%s\nGitCommitLog=%s\nGitStatus=%s\nBuildTime=%s\nGoVersion=%s\nruntime=%s/%s\n",
		Version, GitCommitLog, GitStatus, BuildTime, BuildGoVersion, runtime.GOOS, runtime.GOARCH)
}

// 对一些值做美化处理
func beauty() {
	if GitStatus == "" {
		// GitStatus 为空时，说明本地源码与最近的 commit 记录一致，无修改
		// 为它赋一个特殊值
		GitStatus = "cleanly"
	} else {
		// 将多行结果合并为一行
		GitStatus = strings.Replace(strings.Replace(GitStatus, "\r\n", " |", -1), "\n", " |", -1)
	}
}

func init() {
	beauty()
}
