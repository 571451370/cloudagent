package qga

import "time"

var (
	l      *Logger
	uptime time.Time
)

func init() {
	uptime = time.Now()
}

func GetUptime() time.Time {
	return uptime
}

// FileSystem struct
type FileSystem struct {
	Device  string
	Path    string
	Type    string
	Options []string
}

// ExecStatus struct
type ExecStatus struct {
	Exited   bool   `json:"exited"`
	ExitCode *int   `json:"exitcode,omitempty"`
	Signal   int    `json:"signal,omitempty"`
	OutData  string `json:"out-data,omitempty"`
	ErrData  string `json:"err-data,omitempty"`
	OutTrunc bool   `json:"out-truncated,omitempty"`
	ErrTrunc bool   `json:"err-truncated,omitempty"`
}

const (
	// MaxMessageLength is the maximum message length
	MaxMessageLength = 4 * 1024
)
