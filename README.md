# Go Performance Related Test Scripts
## sync.Pool
* Use `sync.Pool` to reuse objects
* `Pool.Get` will create a new object if pool is empty, and pool.New is non-nil
* How is this different from using global variables?
    * Because you don't know which object is used for how long and which is free.
    * The objects can be Garbage Collected
    * Concurrent-safe and optimized for multiple goroutines

## GOGC, GOMEMLIMIT and GODEBUG
* Try running `GOMAXPROCS=1 GOMEMLIMIT=1GiB GOGC=off GODEBUG=gctrace=1 go run ./gc/main.go`
* Example GC `gc 2 @0.287s 1%: 0.009+30+0.003 ms clock, 0.009+0.43/3.5/0+0.003 ms cpu, 873->1021->196 MB, 922 MB goal, 0 MB stacks, 0 MB globals, 1 P`
  * Third GC run (gc 2)
  * GC was triggered 287 millis after program started
  * GC used 1% of the total available CPU time
  * `0.009+30+003 ms` clock - STW start 0.009 ms, 30 ms concurrent mark and sweep, and 0.003 ms STW end event
  * `0.009+0.43/3.5/0+0.003` ms cpu
    * 0.009 ms — STW start
    * 0.43/3.5/0 ms — concurrent mark/sweep:
      * 0.43 ms in mark assist (mutator doing GC)
      * 3.5 ms in dedicated GC workers
      * 0 ms idle (no unused GC threads)
    * 0.003 ms — STW end
  * `873->1021->196 MB` => Heap Live start, heap grew during GC, and heap live after GC
  * Next target GC size is 922
  * 0 MB stacks, 0 MB globals - Memory in use for goroutine stacks and global variables
  * 1 P - Number of logical processors P

* There will be two GCs pausing for 22 micros (in my machine)
* Running with GOGC=100 means GC running every time heap doubles `GOMAXPROCS=1 GOMEMLIMIT=1GiB GOGC=100 GODEBUG=gctrace=1 go run ./gc/main.go`
  * There will be more frequent GCs even when there is more memory to use
  * This resulted in 24 GC cycles consuming ~ 300 - 600 micros totally
  * I think same theory applies - Do bulk batch processing instead of more frequent small GCs
  * This also tells go to use more available memory. Otherwise it will limit itself to small amount of memory and run GC every time it doubles. Not optimal

* GODEBUG=schedtrace=100 `GOMAXPROCS=12 GOMEMLIMIT=1GiB GOGC=100 GODEBUG=schedtrace=100 go run ./gc/main.go `
  * Example `SCHED 1361ms: gomaxprocs=12 idleprocs=12 threads=23 spinningthreads=0 needspinning=0 idlethreads=18 runqueue=0 [0 0 0 0 0 0 0 0 0 0 0 0]`
  * SCHED 410 - Snapshot of scheduler stats 410 millis after program started
  * idleprocs=0 - Number of idle Ps
  * threads=3 - Total OS threads created by Go runtime
  * runqueue=0 - Number of goroutines waiting to be scheduled globally (not asssigned to any P)
  * [2 0 0 ...] - Local run queue size of P. How many goroutines are ready to run on that P
