# Chapter 3: Go’s concurrency building blocks

## Goroutines
**Goroutine **(def): A function that is running concurrently (remember: not necessarily in parallel!) alongside other code.
  - Every program has at least one: the main goroutine
  - Special class of coroutine

**M:N Scheduler:** maps M green threads to N OS threads. Goroutines are then scheduled on the green threads. This M:N Scheduler is the mechanism for hosting goroutines

- When we have more goroutines than green threads available, the scheduler handles the distribution of the goroutines across the available threads and ensures that when these goroutines become blocked, other goroutines can be run

Go follows a model of concurrency called the *fork-join model*. This is important to grasp as we want our goroutines to finish deterministically.
- Goroutines are not garbage collected
- Goroutines live in the same address space
- Goroutines are cheap, so you should only be discussing their cost if you’ve proven they are the root cause of a performance issue
- goroutine context switching is **MUCH** faster than os thread context switching
  
## The sync Package

### WaitGroup
WaitGroup is a great way to wait for a set of concurrent operations to complete when you either don’t care about the result of the concurrent operation, or you have other means of collecting their results. If neither of those conditions are true, I suggest you use channels and a select statement instead.

**waitgroup** = concurrent-safe counter.
  - calls to `Add` increment
  - calls to `Done` decrement
  - calls to `Wait` block until the counter is zero

Call `Add` outside of coroutine, else its possible that Wait will get called before the go routines are scheduled and it won’t wait at all since nothing was added to the wait group

It’s customary to couple calls to `Add` as closely as possible to the goroutines they’re helping to track, but sometimes you’ll find Add called to track a group of goroutines all at once. I usually do this before for loops like this:
```go        
wg.Add(numGreeters)
for i := 0; i < numGreeters; i++ {
    go hello(&wg, i+1)
}
```
 
### Mutex and RWMutex
You’ll notice that we always call Unlock within a defer statement. This is a very common idiom when utilizing a Mutex to ensure the call always happens, even when panicing. Failing to do so will probably cause your program to deadlock.

```go
increment := func() {
    lock.Lock()                 
    defer lock.Unlock()         
    count++
    fmt.Printf("Incrementing: %d\n", count)
}
```

The sync.RWMutex (sync.Locker) is conceptually the same thing as a Mutex: it guards access to memory; however, RWMutex gives you a little bit more control over the memory.

- You can request a lock for reading, in which case you will be granted access unless the lock is being held for writing. 
- This means that an arbitrary number of readers can hold a reader lock so long as nothing else is holding a writer lock.

### Cond

*"...a rendezvous point for goroutines waiting for or announcing the occurrence of an event."*

Cond type is much more performant than utilizing channels.

```go
c := sync.NewCond(&sync.Mutex{}) 
c.L.Lock() 
for conditionTrue() == false {
    c.Wait() 
}
c.L.Unlock() 
```

There are two methods that the Cond type provides for notifying goroutines blocked on a Wait call that the condition has been triggered. Internally, the runtime maintains a FIFO list of goroutines waiting to be signaled
- **Signal:** finds the goroutine that’s been waiting the longest and notifies that
- **Broadcast:** sends a signal to all goroutines that are waiting

Like most other things in the sync package, usage of Cond works best when constrained to a tight scope, or exposed to a broader scope through a type that encapsulates it.
