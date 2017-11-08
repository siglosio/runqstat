# runqstat
Linux Run Queue Utility

runqstat is a Linux command-line tool written in Go for collecting stats about the CPU, especially the run queue - this essentially gives you the load average without any blocking disk I/O. Note this is an approximation since the kernel does not provide this info directly, so we have to sample the run queue often (every 10ms) and average over 1 or more seconds.

## BETA RELEASE 
Still in testing and features and/or output may change.
In addition, the code is not fully cleaned up.

Usage:

runqstat [-d duration] [-i interval] [-c count] [-q] [-b] [-v] [-h]

  The options are as follows:

       -d      Duration of the sampling run. Default is 1 second.
       -i      Interval of the sample time. Default is 25 milliseconds.
       -c      Count of overall runs. Default is 1 run.
       -q      Get queue only, which will subtract the number of CPU cores. Default is off.
       -b      Include blocked count, from /proc/stats. Default is off.
       -v      Verbose, for debugging and more info. Default is off.
       -h      Help and usage.

## Contributing
We are not ready for contributors until we can get the code cleaned up and standardized for Go best practices.

However, you can contribute by:
- [Report bugs](https://github.com/opsstack/runqstat/issues/new)
- [Improve documentation](https://github.com/opsstack/runqstat/issues?q=is%3Aopen+label%3Adocumentation)
- [Review code and feature proposals](https://github.com/opsstack/runqstat/pulls)

## Installation:

You can download the binaries directly from the [binaries](https://github.com/opsstack/runqstat/binaries) section.  We'll have RPMs and DEB packages as soon as things stabilize a bit.

### From Source:

This is a single source file project for now, so you can just compile as you would any Golang project.

There is a single external dependency, [pflag](https://github.com/ogier/pflag)

## How to use it:

See usage with:

```
./runqstat --help
```
