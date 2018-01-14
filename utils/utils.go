package utils

import(
	"encoding/json"
	"strconv"
)

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

type CPUStats struct{
	TotalRunning int64
	TotalIdle int64
	Total int64
	Performance float64
}

func (s *CPUStats)ToJson() string{
	b,er:=json.MarshalIndent(s," ","\t")
	if(er!=nil){
		return ""
	}
	return string(b)
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

	totalProcessing:=cpuInfoList[0].UserMode+cpuInfoList[0].NicedProcesses+cpuInfoList[0].SystemProcesses+cpuInfoList[0].IrqProcesses+cpuInfoList[0].SoftIrq+cpuInfoList[0].Steal

	totalIdle:=cpuInfoList[0].IdleProcesses+cpuInfoList[0].IowaitProcesses
	total:=totalProcessing+totalIdle

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
