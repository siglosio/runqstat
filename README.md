# runqstat
Linux Run Queue Utility

runqstat is a Linux command-line tool written iin Go for collecting stats about the CPU, especially the run queue

## ALPHA RELEASE 
Not all features fully implemented nor code fully cleaned up.

Usage:

runqstat —-interval —-duration —-blocked —-queue —-method --help 

-i - Interval - Sample interval in milliseconds, default 10ms
-d - Duration - How long to run for, default 1 second
-m - Method - Default is Average, can also do wma Weighted Moving Average. Weights are fixed.
-b - Blocked - Show second line with blocked count
-q - Queue Only - Show ‘run queue’ - 'core count’ to get real queue, as default is to show the running & queued.
-v - Verbose - Show additional info for testing and internals
-h - Help - Usage help


## Contributing
We are not ready for contributors until we can get the code cleaned up and standardized or Go best practices.

However, you can contribute by:
- Fix and [report bugs](https://github.com/opsstack/runqstat/issues/new)
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
