package response

import (
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

type ProcessListP struct {
	ProcessList []ProcessInfoP `json:"process_list"`
}

type ProcessInfoP struct {
	Pid           int32                   `json:"pid"`
	PPid          int32                   `json:"ppid"`
	Name          string                  `json:"name"`
	UserName      string                  `json:"user_name"`
	Status        string                  `json:"status"`
	CpuPercent    float64                 `json:"cpu_percent"`
	MemoryPercent float32                 `json:"memory_percent"`
	MemoryUse     uint64                  `json:"memory_use"`
	IOCounters    *process.IOCountersStat `json:"io_counters"`
	NumThreads    int32                   `json:"num_threads"`
	NumFDs        int32                   `json:"num_fds"`
}

type SystemServiceInfo struct {
	Name string `json:"name"`
	R0   bool   `json:"r0"`
	R1   bool   `json:"r1"`
	R2   bool   `json:"r2"`
	R3   bool   `json:"r3"`
	R4   bool   `json:"r4"`
	R5   bool   `json:"r5"`
	R6   bool   `json:"r6"`
	Ps   string `json:"ps"`
}

type ConnectionListP struct {
	List       []*ConnectionInfo `json:"list"`
	Statistics []net.IOCountersStat
}

type ConnectionInfo struct {
	net.ConnectionStat
	ProcessName string `json:"process_name"`
}

type LinuxUserListP struct {
	List []*LinuxUserInfo `json:"list"`
}

type LinuxUserInfo struct {
	Username string `json:"username"`
	Home     string `json:"home"`
	Group    string `json:"group"`
	UID      int    `json:"uid"`
	GID      int    `json:"gid"`
	Shell    string `json:"shell"`
	Ps       string `json:"ps"`
}
