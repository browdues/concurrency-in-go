# Chapter 1. An Introduction to Concurrency

## Moore’s Law, Web Scale, and the Mess We’re In

Amdahl's law: gains are bounded by how much of the program must be written in a sequential manner.

Web Scale Software: Able to handle hundreds of thousands (or more) of simultaneous workloads by adding more instances of the application. (*Embarrasingly Parallel*)

## Why Is Concurrency Hard?

### Race Conditions - 2+ operations must exec in correct order, but program isn't written to guarantee this order is maintained. (data race) 

When writing concurrent code
- Meticulously iterate through possible scenarios
- *HelpfulTip:* Imagine an hour passes between the time when the goroutine is invoked, and when it is run
- Target **Logical Correctness**, not hacks that 'make it work'

### Atomicity - *(property of)* - within the context that it is operating, it is indivisible, or uninterruptible.

**keyword: context**

- “indivisible” and “uninterruptible.” These terms mean that within the context you’ve defined, something that is atomic will happen in its entirety without anything happening in that context simultaneously
- important because if something is atomic, implicitly it is safe within concurrent contexts.

Define the context or scope when thinking about atomicity, because the atomicity of an operation can change depending on the currently defined scope.
-> When thinking about atomicity, very often the first thing you need to do is to define the context, or scope, the operation will be considered to be atomic in. 


### Memory Access Synchronization

Critical Section: The part of your program that needs exclusive access to a shared resource.

In order to get this right, need to answer 2 questions:

1. Are my critical sections entered and exited repeatedly?
2. What size should my critical sections be?

### Deadlocks, Livelocks, and Starvation
Prev. sections about correctness. Even if this is achieved, there is another class of issues...

**DEADLOCK** A deadlocked program is one in which all concurrent processes are waiting on one another

*Coffman Condtions:* conditions that must be present for deadlocks

1. Mutual Exclusion: a conc proc holds exclusive rights to a resource at any time
2. Wait For Condition: A conc proc mus simultaneously hold a resource and be waiting for an additional resource.
3. No Preemption: A resource held by a concurrent process can only be released by that process so it fulfills this condition.
4. Circular Wait: A concurrent process (P1) must be waiting on a chain of other concurrent processes (P2), which are in turn waiting on it (P1), so it fulfills this final condition too.

If we ensure that at least one of these conditions is **not** true, we can prevent deadlocks from occuring.

**LIVELOCK**
That moment when you both move out of the way, but you're in the way of each other, and this dance never ends

These can be harder to spot than deadlocks. Since CPU is active, it appears that work is occuring. Livelocks are a subset of a larger set of problems known as Starvation.

**STARVATION** Any situation where a concurrent process can't get all the resources it needs to perform work.

Referencing the example:  Metrics are a good way for detecting and solving starvation. Log when work is accomplished and then determine if the rate of work is as high as expected.

### Determining Concurrency Safety
Humans are involved, so be careful to communicate. For concurrent functions, make sure comments cover the following

mportantly, the comment covers these aspects:

1. Who is responsible for the concurrency?
2. How is the problem space mapped onto concurrency primitives?
3. Who is responsible for the synchronization?

## Simplicity in the Face of Complexity
