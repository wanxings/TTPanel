package system

import (
	"TTPanel/internal/global"
	"TTPanel/internal/model/response"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"time"
)

func GetIoStat() map[string]disk.IOCountersStat {
	var AllIOCounter disk.IOCountersStat
	NewIOCounter := make(map[string]disk.IOCountersStat)
	cacheTimeout := 1 * time.Minute
	mtime := time.Now().Unix()

	ioKey := "iostat"
	ioTimeKey := "iostatTime"

	diskIo, _ := global.GoCache.Get(ioKey)
	diskIoTime, _ := global.GoCache.Get(ioTimeKey)

	oldDiskIo, ok := (diskIo).(map[string]disk.IOCountersStat)
	if ok == false {
		diskIoTime = mtime
		oldDiskIo = nil
	}
	diskIoTimeS, _ := diskIoTime.(int64)
	sTime := uint64(mtime - diskIoTimeS)
	if sTime == 0 {
		sTime = 1
	}

	IOCounters, _ := disk.IOCounters()
	if oldDiskIo == nil {
		oldDiskIo = IOCounters
	}
	for _, IOCounter := range IOCounters {
		item := disk.IOCountersStat{
			ReadCount:        (IOCounter.ReadCount - oldDiskIo[IOCounter.Name].ReadCount) / sTime,
			ReadBytes:        (IOCounter.ReadBytes - oldDiskIo[IOCounter.Name].ReadBytes) / sTime,
			ReadTime:         (IOCounter.ReadTime - oldDiskIo[IOCounter.Name].ReadTime) / sTime,
			WriteCount:       (IOCounter.WriteCount - oldDiskIo[IOCounter.Name].WriteCount) / sTime,
			WriteBytes:       (IOCounter.WriteBytes - oldDiskIo[IOCounter.Name].WriteBytes) / sTime,
			WriteTime:        (IOCounter.WriteTime - oldDiskIo[IOCounter.Name].WriteTime) / sTime,
			IoTime:           (IOCounter.IoTime - oldDiskIo[IOCounter.Name].IoTime) / sTime,
			MergedReadCount:  (IOCounter.MergedReadCount - oldDiskIo[IOCounter.Name].MergedReadCount) / sTime,
			MergedWriteCount: (IOCounter.MergedWriteCount - oldDiskIo[IOCounter.Name].MergedWriteCount) / sTime,
		}
		NewIOCounter[IOCounter.Name] = item

		AllIOCounter.ReadCount = AllIOCounter.ReadCount + item.ReadCount
		AllIOCounter.MergedReadCount = AllIOCounter.MergedReadCount + item.MergedReadCount
		AllIOCounter.WriteCount = AllIOCounter.WriteCount + item.WriteCount
		AllIOCounter.MergedWriteCount = AllIOCounter.MergedWriteCount + item.MergedWriteCount
		AllIOCounter.ReadBytes = AllIOCounter.ReadBytes + item.ReadBytes
		AllIOCounter.WriteBytes = AllIOCounter.WriteBytes + item.WriteBytes
		AllIOCounter.ReadTime = AllIOCounter.ReadTime + item.ReadTime
		AllIOCounter.WriteTime = AllIOCounter.WriteTime + item.WriteTime
		AllIOCounter.IoTime = AllIOCounter.IoTime + item.IoTime
	}
	NewIOCounter["all"] = AllIOCounter

	global.GoCache.Set(ioKey, IOCounters, cacheTimeout)
	global.GoCache.Set(ioTimeKey, mtime, cacheTimeout)

	return NewIOCounter
}

func GetCPU() (*response.CpuInfo, *[]cpu.TimesStat) {
	cpuInfos, _ := cpu.Info()
	cpuCoresCount, _ := cpu.Counts(false)
	cpuLogicalCount, _ := cpu.Counts(true)
	cpuPercent, _ := cpu.Percent(0, true)
	cpuTimes, _ := cpu.Times(false)
	return &response.CpuInfo{
		Name:         cpuInfos[0].ModelName,
		CoresCount:   cpuCoresCount,
		LogicalCount: cpuLogicalCount,
		Percents:     cpuPercent,
	}, &cpuTimes
}

func GetMemory() *response.Memory {
	memoryInfo, _ := mem.VirtualMemory()
	return &response.Memory{
		Total:   memoryInfo.Total,
		Buffers: memoryInfo.Buffers,
		Cached:  memoryInfo.Cached,
		Used:    memoryInfo.Used,
		Free:    memoryInfo.Free,
		Percent: memoryInfo.UsedPercent,
	}
}

func GetNetWork() map[string]*response.NetStat {
	cacheTimeout := 10 * time.Minute
	oTime, _ := global.GoCache.Get("oTime")
	nTime := time.Now().Unix()
	NetWorkInfos, _ := net.IOCounters(true)
	NetWorkData := make(map[string]*response.NetStat)
	var AllStat response.NetStat

	for _, NetWorkInfo := range NetWorkInfos {

		upKey := fmt.Sprintf("%v_up", NetWorkInfo.Name)
		downKey := fmt.Sprintf("%v_down", NetWorkInfo.Name)
		oTimeKey := "oTime"

		if oTime == nil {
			oTime = nTime
			global.GoCache.Set(upKey, NetWorkInfo.BytesSent, cacheTimeout)
			global.GoCache.Set(downKey, NetWorkInfo.BytesRecv, cacheTimeout)
			global.GoCache.Set(oTimeKey, oTime, cacheTimeout)
		}
		oTimeS, _ := oTime.(int64)
		up, _ := global.GoCache.Get(upKey)
		down, _ := global.GoCache.Get(downKey)
		if up == nil {
			up = NetWorkInfo.BytesSent
		}
		upS, _ := up.(uint64)
		if down == nil {
			down = NetWorkInfo.BytesRecv
		}
		downS, _ := down.(uint64)

		sTime := uint64(nTime - oTimeS)
		if sTime == 0 {
			sTime = 1
		}
		NetWorkData[NetWorkInfo.Name] = &response.NetStat{
			Down:        (NetWorkInfo.BytesRecv - downS) / sTime,
			DownPackets: NetWorkInfo.PacketsSent,
			DownTotal:   NetWorkInfo.BytesRecv,
			Up:          (NetWorkInfo.BytesSent - upS) / sTime,
			UpPackets:   NetWorkInfo.PacketsRecv,
			UpTotal:     NetWorkInfo.BytesSent,
		}
		AllStat.Down = AllStat.Down + NetWorkData[NetWorkInfo.Name].Down
		AllStat.DownPackets = AllStat.DownPackets + NetWorkData[NetWorkInfo.Name].DownPackets
		AllStat.DownTotal = AllStat.DownTotal + NetWorkData[NetWorkInfo.Name].DownTotal
		AllStat.Up = AllStat.Up + NetWorkData[NetWorkInfo.Name].Up
		AllStat.UpPackets = AllStat.UpPackets + NetWorkData[NetWorkInfo.Name].UpPackets
		AllStat.UpTotal = AllStat.UpTotal + NetWorkData[NetWorkInfo.Name].UpTotal

		global.GoCache.Set(upKey, NetWorkInfo.BytesSent, cacheTimeout)
		global.GoCache.Set(downKey, NetWorkInfo.BytesRecv, cacheTimeout)
		global.GoCache.Set(oTimeKey, time.Now().Unix(), cacheTimeout)

	}
	NetWorkData["all"] = &AllStat
	return NetWorkData
}

// GetLoad 获取负载
func GetLoad() *load.AvgStat {
	loadAvgInfo, _ := load.Avg()
	return loadAvgInfo
}

func GetDisk() *[]response.Disk {
	var disks []response.Disk
	diskInfos, _ := disk.Partitions(false)

	//diskIOCounters, _ := disk.IOCounters()
	for _, diskInfo := range diskInfos {
		diskUseInfo, _ := disk.Usage(diskInfo.Mountpoint)
		disks = append(disks, response.Disk{
			Device:     diskInfo.Device,
			Type:       diskInfo.Fstype,
			MountPoint: diskInfo.Mountpoint,
			InodeUsage: &response.DiskUsage{
				Total:   diskUseInfo.InodesTotal,
				Used:    diskUseInfo.InodesUsed,
				Free:    diskUseInfo.InodesFree,
				Percent: diskUseInfo.InodesUsedPercent,
			},
			SpaceUsage: &response.DiskUsage{
				Total:   diskUseInfo.Total,
				Used:    diskUseInfo.Used,
				Free:    diskUseInfo.Free,
				Percent: diskUseInfo.UsedPercent,
			},
		})
	}
	return &disks
}
func GetSystemInfo() string {
	hostInfo, _ := host.Info()
	return fmt.Sprintf("%v %v %v %v", hostInfo.OS, hostInfo.Platform, hostInfo.PlatformVersion, hostInfo.KernelArch)
}
func GetSwapMemory() *mem.SwapMemoryStat {
	data, _ := mem.SwapMemory()
	return data
}
