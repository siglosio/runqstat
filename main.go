// Copyright 2017-2018 OpsStack.io. All rights reserved.
// Use of this source code is governed by the GPL
// license that can be found in the LICENSE file.

// TODO
//   Error handling for nice output

package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"strings"
	"time"
	"flag"
	"runtime"
	"bufio"
	"encoding/json"
	"strconv"
)

// constants
const (
	version = "0.2"
	defaultDuration = 1 // seconds
	defaultInterval = 10 // milliseconds
	defaultMethod = "average"
	defaultBlocked = false
	procStatPath = "/proc/stat"
)

// CLI flags
var (
	flagDuration int64 = 0
	flagInterval int64 = 0
	flagMethod string = "average"
	flagBlocked bool = false
	flagVerbose	bool = false
	duration int64
	interval int64
	CpuCount int
	cpuInfoList []CPUInfo
)

func main() {

	var (
		loopCount uint64 // Times we need to loop
	)
	cpuInfoList = make([] CPUInfo, CpuCount)

	var totalRunning float64	= 0.0
	var totalBlocked uint64 	= 0
	var running float64 		= 0.0
	var blocked uint64 			= 0
	var i uint64 				= 0
	var average float64 		= 0.0

	fmt.Printf("runq Version %s - Copyright 2017-2018 OpsStack Inc.\n", version)

	// Calc how many times we'll loop for our duration & interval
	loopCount = uint64(float64(duration) / (float64(interval) / 1000))

	// Master loop to gather data
	for i = 0; i < loopCount; i++ {
		running, blocked = getCPUSample()
		time.Sleep(defaultInterval * time.Millisecond)
		totalRunning = totalRunning + running
		totalBlocked = totalBlocked + blocked
	}

	if i > 0 {
		fmt.Printf("Loops: %d\n", i)
		average = float64(totalRunning / float64(i))
	} else {
		average = 0
	}

	fmt.Printf("CPU usage average: %f\n", average)
}

// Get our CPU samples from /proc
func getCPUSample() (running float64, blocked uint64) {

	running = 0.0
	contents, err := ioutil.ReadFile(procStatPath)
	if err != nil {
		return
	}
	reader:=bufio.NewScanner(strings.NewReader(string(contents)));
	cpus := 0
	var t[] string

	// Read file, one line per CPU up to cpu count
	for reader.Scan() && cpus < CpuCount{
		text := reader.Text()
		t = strings.Fields(text)
		cpuInfoList[cpus].updateStats(t)
		cpus++
	}

	fmt.Printf("Performance: %v\n", cpuInfoList[0].ToJson())
	running=cpuInfoList[0].Stats.Performance

	return
}

/*
*Initialize application statistics
*/
func init() {
	// Setup arguments, must do before calling Parse()
	flag.Int64Var(&flagDuration, "duration", 0, "Duration (s)")
	flag.Int64Var(&flagInterval, "interval", 0,"Interval (ms)")
	flag.BoolVar(&flagBlocked, "blocked", false, "Include Blocked")
	flag.StringVar(&flagMethod, "method", "average", "Method")
	flag.BoolVar(&flagVerbose, "verbose", false, "Verbose")

	flag.Parse() // Process argurments

	if flagDuration > 0 {
		duration = flagDuration
		fmt.Printf("Duration: %d \n", duration)
	} else {
		duration = defaultDuration
	}

	if flagInterval > 0 {
		interval = flagInterval
		fmt.Printf("Interval: %d \n", interval)
	} else {
		interval = defaultInterval
	}

	CpuCount=runtime.NumCPU() // Count CPU cores, needed later for stats
	fmt.Printf("====== Number of CPUs: %d ============\n",CpuCount)
}

func argsCheck() {
	// if user does not supply flags, print usage
	if flag.NFlag() == 0 {
		fmt.Printf("Usage: %s [options]\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

/**
*Update stats for a given CPU
*Ref: https://stackoverflow.com/questions/23367857/accurate-calculation-of-cpu-usage-given-in-percentage-in-linux
*/

func (c *CPUInfo) updateStats(t []string ){
	c.Count++
	c.Name 				= t[0]
	c.UserMode 			= toInt64(t[1])
	c.NicedProcesses 	= toInt64(t[2])
	c.SystemProcesses 	= toInt64(t[3])
	c.IdleProcesses 	= toInt64(t[4])
	c.IowaitProcesses 	= toInt64(t[5])
	c.IrqProcesses 		= toInt64(t[6])
	c.SoftIrq 			= toInt64(t[7])
	c.Steal 			= toInt64(t[8])

	totalProcessing :=
		cpuInfoList[0].UserMode +
		cpuInfoList[0].NicedProcesses +
		cpuInfoList[0].SystemProcesses +
		cpuInfoList[0].IrqProcesses +
		cpuInfoList[0].SoftIrq +
		cpuInfoList[0].Steal

	totalIdle :=
		cpuInfoList[0].IdleProcesses +
		cpuInfoList[0].IowaitProcesses

	total := totalProcessing + totalIdle

	//Skip first run
	if(c.Count>1){
		c.Stats.Performance=float64((total-c.Stats.Total)-(totalIdle-c.Stats.TotalIdle))/float64(total-c.Stats.Total)
	}

	c.Stats.Total 			= total
	c.Stats.TotalRunning 	= totalProcessing
	c.Stats.TotalIdle 		= totalIdle
}

func (c *CPUInfo) ToJson() string{
	b,er:=json.MarshalIndent(c," ","\t")
	if(er!=nil){
		return ""
	}
	return string(b)
}

/**
*Convert string to int64
*/
func toInt64(t string) int64{
	val,err:=strconv.ParseInt(t,10,64)
	if(err!=nil){
		return 0
	}
	return val
}

func (s *CPUStats)ToJson() string{
	b,er:=json.MarshalIndent(s," ","\t")
	if(er!=nil){
		return ""
	}
	return string(b)
}

type CPUStats struct{
	TotalRunning int64
	TotalIdle int64
	Total int64
	Performance float64
}

type CPUInfo struct{
	Name string
	UserMode int64 //Executing in user mode
	NicedProcesses int64  //Niced Processes executing in user mode
	SystemProcesses int64//System processed executing in kernel mode
	IdleProcesses int64
	IowaitProcesses int64
	IrqProcesses int64
	SoftIrq int64
	Steal int64
	GuestMode int64
	Stats CPUStats
	Count int64
}
