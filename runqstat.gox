// Copyright 2017 OpsStack.io. All rights reserved.
// Use of this source code is governed by the GPL
// license that can be found in the LICENSE file.

// TODO
//   Error handling for nice output

package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
	flag "github.com/ogier/pflag"
	"runtime"
)

// constants
const (
	version              = "0.1"
	procStatPath         = "/proc/stat"
	procSchedStatPath    = "/proc/schedstat"
	procStatFieldRunning = "procs_running"
	procStatFieldBlocked = "procs_blocked"
	procCPUPrefix        = "cpu"
	copyright            = "Copyright 2017 by OpsStack.io"
)

// Global vars
var (
	argDuration int64
	argInterval int64
	argCount    int64
	argMethod   string
	flagQueue   bool
	flagLatency bool
	flagBlocked bool
	flagVerbose bool
	flagHelp    bool
)

////////// Program Main ////////////

func main() {

	var (
		innerLoopCount       int64 // Times we need to loop inside
		totalRunning         int64 = 0
		totalBlocked         int64 = 0
		running              int64 = 0
		blocked              int64 = 0
		inner                int64
		outer                int64
		averageRunning       float64
		averageQueue         float64
		averageBlocked       float64
		cpuCount             float64
		lastLatency          int64
		newLatency           int64
		deltaLatency         int64
		averageLatency       int64
		averageLatencyOutput string
	)

	// Check our command-line arguments
	argsCheck(version, copyright);

	// Calculate how many loops we need for our interval
	innerLoopCount = int64(float64(argDuration) / float64(argInterval) * 1000)

	cpuCount = float64(runtime.NumCPU())

	if flagVerbose {
		fmt.Printf("Running %d count for %d second(s) with sample interval %dms\n", argCount, argDuration, argInterval)
		fmt.Printf("Inner Loop Count %d with Duration %ds and Interval %dms\n", innerLoopCount, argDuration, argInterval)
		fmt.Printf("CPU Count: %d\n", int64(cpuCount))
		fmt.Println()
	}

	// Print headers if count > 1
	if argCount > 1 {
		fmt.Print("runq_avg")
		if flagBlocked {
			fmt.Print(" runq_blocked")
		}
		if flagLatency {
			fmt.Print(" runq_latency")
		}
		fmt.Println()
	}

	// Outer loop for counts, repeats each duration
	lastLatency = 0
	for outer = 0; outer < argCount; outer++ {

		totalRunning = 0
		totalBlocked = 0

		for inner = 0; inner < innerLoopCount; inner++ {
			running, blocked = getCPURunning()
			totalRunning = totalRunning + running
			totalBlocked = totalBlocked + blocked
			// Main sleep timer
			time.Sleep(time.Duration(argInterval) * time.Millisecond)
		}

		if inner > 0 {
			averageRunning = float64(totalRunning / inner)
			averageBlocked = float64(totalBlocked / inner)
		} else {
			averageRunning = 0
			averageBlocked = 0
		}

		// Subtract CPU cores to get the real queue
		if averageRunning > cpuCount {
			averageQueue = averageRunning - cpuCount
		} else {
			averageQueue = 0
		}

		// Process Latency
		if flagLatency {
			newLatency = getCPULatency()

			deltaLatency = newLatency - lastLatency
			averageLatency = int64(float64(deltaLatency) / (float64(argDuration)) * 1000) // Milliseconds
			//fmt.Printf("new, last, deltal, averageLatency: %d - %d - %d - %d\n", newLatency, lastLatency, deltaLatency, averageLatency)
			lastLatency = newLatency

			// Special processing for first pass as we dont' have previous sample
			// Blank out first response
			if outer == 0 {
				averageLatencyOutput = "-"
			} else {
				averageLatencyOutput = fmt.Sprintf("%d", averageLatency)
			}
		}

		// Print our output
		if flagQueue {
			fmt.Printf("%6.2f", averageQueue)
		} else {
			fmt.Printf("%6.2f", averageRunning)

		}
		if flagBlocked {
			fmt.Printf("     %6.2f", averageBlocked)
		}
		if flagLatency {
			fmt.Printf("     %6s", averageLatencyOutput)
		}
		fmt.Println()

	} // Outer loop

} // Main()

// Process arguments
func argsCheck(version string, copyright string) {
	if flagHelp {
		fmt.Printf("runqstat Version %s - %s\n\n", version, copyright)
		fmt.Printf("Usage: %s [options]\n\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

// Inner sampling loop to read instantaneous values
func getCPURunning() (running, blocked int64) {

	var lines []string
	running = 0
	blocked = 0

	contents, err := ioutil.ReadFile(procStatPath)

	if err != nil {
		//		fmt.Println("Can't find /proc/stat")
		//return
		//		lines = make([]string, "procs_running 33")
		// Test data for OSX
		lines = []string{"procs_running 8"}
	} else {
		lines = strings.Split(string(contents), "\n")
	}

	// Parsing - could be improved
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == procStatFieldRunning {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				//fmt.Println("I, Val %d %d", i, val) //Debug our parser
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				if i == 1 { // get 2nd field
					running = running + int64(val)
				}
			}

			if fields[0] == procStatFieldBlocked && flagBlocked {
				numFields := len(fields)
				for i := 1; i < numFields; i++ {
					val, err := strconv.ParseUint(fields[i], 10, 64)
					//fmt.Println("I, Val %d %d", i, val) //Debug our parser
					if err != nil {
						fmt.Println("Error: ", i, fields[i], err)
					}
					if i == 1 { // get 2nd field
						blocked = blocked + int64(val)
					}
				}
			}

			//fmt.Println("Running %d, Blocked %d", running, blocked) //Debug our parser
			return
		}
	}
	return
}

// Parse and sum CPUs in /proc/schedstats to get task latency
func getCPULatency() (sumLatency int64) {

	var (
		lines []string
	)

	contents, err := ioutil.ReadFile(procSchedStatPath)

	if err != nil {
		//		fmt.Println("Can't find /proc/stat")
		//return
		// Test data for OSX
		lines = []string{
			"version 15",
			"timestamp 19218408415",
			"cpu0 379022 0 4137731047 1554123345 2434574576 2434574576 140081641042736 62760166051387 2583041356",
			"random",
			"cpu1 379022 0 4137731047 1554123345 2434574576 2434574576 140081641042736 62760166051387 2583041356",
			""}
	} else {
		lines = strings.Split(string(contents), "\n")
	}

	// Parsing - could be improved
	sumLatency = 0
	for _, line := range lines {
		fields := strings.Fields(line)
		//fmt.Printf("Line: %s\n", line)
		// Need to check length or else will get index error on blank line
		if len(fields) > 0 && strings.Contains(fields[0], procCPUPrefix) {
			//fmt.Printf("fields[0] %s\n", fields[0])
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				//fmt.Printf("fields[%d], Value: %d\n", i, val) //Debug our parser
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				if i == 7 { // get 8th field (7th data field, after cpuX)
					//fmt.Printf("Latency field: %d\n", val)
					sumLatency = sumLatency + int64(val)
				}
			}
		}
	}

	// Convert Jiffies to milliseonds (1 Jiffie = 10ms)
	sumLatency = sumLatency * 10

	return
}

// init is called automatically at start
func init() {

	// Setup arguments, must do before calling Parse()
	flag.Int64VarP(&argDuration, "duration", "d", 1, "Duration (s)")
	flag.Int64VarP(&argInterval, "interval", "i", 10, "Interval (ms)")
	flag.Int64VarP(&argCount, "count", "c", 1, "Loop Count")
	flag.BoolVarP(&flagQueue, "queue", "q", false, "Queue Only")
	flag.BoolVarP(&flagLatency, "latency", "l", false, "Latency")
	flag.BoolVarP(&flagBlocked, "blocked", "b", false, "Include Blocked")
	flag.BoolVarP(&flagVerbose, "verbose", "v", false, "Verbose Output")
	flag.BoolVarP(&flagHelp, "help", "h", false, "Help")

	// Future
	//flag.StringVarP(&argMethod, "method", "m", "average", "Averaging Method")

	flag.Parse() // Process argurments
}
