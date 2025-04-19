# Go Performance Related Test Scripts
* Use `sync.Pool` to reuse objects
* `Pool.Get` will create a new object if pool is empty, and pool.New is non-nil
* How is this different from using global variables?
    * Because you don't know which object is used for how long and which is free.
    * The objects can be Garbage Collected
    * Concurrent-safe and optimized for multiple goroutines
