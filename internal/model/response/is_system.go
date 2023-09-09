package response

type CpuInfo struct {
	Name         string    `json:"name"`          //cpu名称
	CoresCount   int       `json:"cores_count"`   //cpu物理核心
	LogicalCount int       `json:"logical_count"` //cpu逻辑核心
	Percents     []float64 `json:"percent"`       //每个核心的使用率
}
type Memory struct {
	Total   uint64  `json:"total"`   //总内存
	Buffers uint64  `json:"buffers"` //缓冲区
	Cached  uint64  `json:"cached"`  //已缓存
	Used    uint64  `json:"used"`    //已使用
	Free    uint64  `json:"free"`    //可用
	Percent float64 `json:"percent"` //使用率
}
type Disk struct {
	Device     string     `json:"device"`      //文件系统
	Type       string     `json:"type"`        //类型
	MountPoint string     `json:"mount_point"` //挂载点
	InodeUsage *DiskUsage `json:"inode_usage"` //磁盘Inode信息
	SpaceUsage *DiskUsage `json:"space_usage"` //磁盘空间信息
}
type DiskUsage struct {
	Total   uint64  `json:"total"`   //容量
	Used    uint64  `json:"used"`    //已使用
	Free    uint64  `json:"free"`    //可用
	Percent float64 `json:"percent"` //使用率
}
type NetStat struct {
	Down        uint64 `json:"down"`         //容量
	DownPackets uint64 `json:"down_packets"` //容量
	DownTotal   uint64 `json:"down_total"`   //容量
	Up          uint64 `json:"up"`           //容量
	UpPackets   uint64 `json:"up_packets"`   //容量
	UpTotal     uint64 `json:"up_total"`     //容量
}
