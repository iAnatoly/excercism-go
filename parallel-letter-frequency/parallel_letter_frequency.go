/*

There are 3 different approaches implemnted below:
* Concurrency with a mutex and a waitgroup
* Concurrency with an array of channels
* Concurrency with a single channel

Benchmarks are:

BenchmarkSequentialFrequency-16                	     261	   4483180 ns/op	    4175 B/op	      14 allocs/op
BenchmarkMutexConcurrentLarge-16               	      61	  19880896 ns/op	    4347 B/op	      16 allocs/op
BenchmarkMultipleChannelsConcurrentLarge-16    	     846	   1350947 ns/op	  227754 B/op	    2328 allocs/op
BenchmarkConcurrentLarge-16                    	    1170	   1027587 ns/op	  123132 B/op	    1238 allocs/op
*/

package letter

import (
	"runtime"
	"sync"
)

// FreqMap records the frequency of each rune in a given text.
type FreqMap map[rune]int

// Frequency counts the frequency of each rune in a given text and returns this
// data as a FreqMap.
func Frequency(s string) FreqMap {
	m := FreqMap{}
	for _, r := range s {
		m[r]++
	}
	return m
}

// -------------------- approach 1:  mutex + waitgroup ----------------------
// safeMap is a dataobject to keep mutex, waitgroup and hashmap together
type safeMap struct {
	freqmap FreqMap
	mutex   sync.Mutex
	wait    sync.WaitGroup
}

// safeCount implements thread-safe counting goroutine
// It uses a mutex field to achieve a lock on the map
// and waitgroup field to report the completion back to the master thread
func safeCount(s string, sMap *safeMap) {
	defer sMap.wait.Done()
	for _, r := range s {
		sMap.mutex.Lock()
		sMap.freqmap[r]++
		sMap.mutex.Unlock()
	}
}

// MutexConcurrentFrequency counts the letter frequency on an array of strings
// Each elemet in the arry of strings is processed concurrently
// MutexConcurrentFrequency returns a map of letters (runes) to their respective frequency
func MutexConcurrentFrequency(ss []string) FreqMap {
	smap := safeMap{freqmap: make(FreqMap)}

	for _, s := range ss {
		smap.wait.Add(1)
		go safeCount(s, &smap)
	}
	smap.wait.Wait() // wait for all SafeCount threads to finish
	return smap.freqmap
}

// ------------------ approach 2 - multiple channels and merge ---------------------
// MultipleChannelsConcurrentFrequency is a function implementing conncurrent frequency calculation.
// The results are sent to the paret function via a channel (channel per goroutine)
func MultipleChannelsConcurrentFrequency(ss []string) FreqMap {
	result := FreqMap{}
	channels := make([]chan FreqMap, len(ss))

	for i := 0; i < len(ss); i++ {
		channels[i] = make(chan FreqMap)
		go func(s string, channel chan FreqMap) {
			channel <- Frequency(s)
		}(ss[i], channels[i])
	}
	for _, m := range channels {
		for k, v := range <-m {
			result[k] += v
		}
	}

	return result
}

// ------------------ approach 3 - single channel and merge ---------------------
// reusing channelCount function from approach #2

// ConcurrentFrequency is a function implementing conncurrent frequency calculation.
// The results are sent to the paret function via a ibuffered channel (shared across all goroutines)
func ConcurrentFrequency(ss []string) FreqMap {
	result := FreqMap{}
	concurrency := runtime.NumCPU()
	if concurrency > 10 {
		concurrency = 10
	}
	channel := make(chan FreqMap, concurrency)

	for _, s := range ss {
		go func(s string) {
			channel <- Frequency(s)
		}(s)
	}
	for range ss {
		for k, v := range <-channel {
			result[k] += v
		}
	}

	return result
}
