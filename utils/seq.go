package utils

import (
	_ "encoding/json"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	_ "github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	_ "github.com/shirou/gopsutil/v3/process"
)

func flow(rMap, tMap map[string]int64) {
	ethx := "all"
	timer := time.NewTicker(time.Second * 5)
	var rLast, tLast int64
	var recv, send int64
	for {
		<-timer.C
		if r, ok := rMap[ethx]; ok {
			fmt.Println(r, rLast)
			if rLast != 0 {
				recv = r - rLast
			}

			rLast = r
		} else {
			fmt.Println("not found")
		}

		if t, ok := tMap[ethx]; ok {
			if tLast != 0 {
				send = t - tLast
			}
			tLast = t
		}
		flow_in := float32(recv / 1024)
		flow_out := float32(send / 1024)
		fmt.Println(flow_in, " ks : ", flow_out, " ks")

	}

}

func Seq() {

	v, _ := mem.VirtualMemory()
	timestamp, _ := host.BootTime()
	t := time.Unix(int64(timestamp), 0)
	platform, family, version, _ := host.PlatformInformation()
	physicalCnt, _ := cpu.Counts(false)
	totalPercent, _ := cpu.Percent(3*time.Second, false)

	// netlist, _ := net.Interfaces()

	// for _, v := range netlist {
	// 	fmt.Println(v.String())
	// }
	incounters, _ := net.IOCounters(false)
	var rMap = make(map[string]int64)
	var tMap = make(map[string]int64)
	for _, v := range incounters {

		if v.Name == "all" {
			rMap["all"] = int64(v.BytesSent)
			tMap["all"] = int64(v.BytesRecv)
			fmt.Println(v)
		}

	}
	flow(rMap, tMap)
	// fmt.Println(incounters)
	fmt.Println("system_time", t.Local().Format("2006-01-02 15:04:05"))
	fmt.Println("server_version:", platform+family+version)
	fmt.Printf("cpu:count:%d \n", physicalCnt)
	fmt.Printf("cpu:total:%v", totalPercent)
	fmt.Printf("mem:total: %v, mem:free:%v, mem:used:%f%%\n", v.Total, v.Free, v.UsedPercent)
	// disk
	// mapStat, _ := disk.IOCounters()
	// for name, stat := range mapStat {
	// 	fmt.Println(name)
	// 	data, _ := json.MarshalIndent(stat, "", "  ")
	// 	fmt.Println(string(data))
	// }

}
