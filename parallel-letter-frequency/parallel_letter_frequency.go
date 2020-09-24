package letter

import (
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

// --- approach 1 - bruteforce
// naiveConcurrentFrequency is a naive implementation w/o concurrency.
// Not used anywhere, created just for reference.
// Curiosly enough, this naive implementation passes the benchmarks.
// BenchmarkConcurrentFrequency-16    	   33374	     35933 ns/op	    2099 B/op	      12 allocs/op
func naiveConcurrentFrequency(ss []string) FreqMap {
	m := FreqMap{}
	for _, s := range ss {
		for _, r := range s {
			m[r]++
		}
	}
	return m
}

// -------------------- approach 2 - mutex + waitgroup ----------------------
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
// BenchmarkConcurrentFrequency-16    	   19760	     61115 ns/op	    2132 B/op	      13 allocs/op
func MutexConcurrentFrequency(ss []string) FreqMap {
	smap := safeMap{freqmap: make(FreqMap)}

	for _, s := range ss {
		smap.wait.Add(1)
		go safeCount(s, &smap)
	}
	smap.wait.Wait() // wait for all SafeCount threads to finish
	return smap.freqmap
}

// ------------------ approach 3 - channels and waitgroups and merge ---------------------
func channelCount(s string, freqChan chan FreqMap) {
	freqmap := FreqMap{}
	for _, r := range s {
		freqmap[r]++
	}
	freqChan <- freqmap
}

// ConcurrentFrequency is a function implementing conncurrent frequency calculation.
// The results are sent to the paret function via a channel
// BenchmarkConcurrentFrequency-16    	   22984	     51336 ns/op	    8328 B/op	      47 allocs/op
func ConcurrentFrequency(ss []string) FreqMap {
	result := FreqMap{}
	channels := make([]chan FreqMap, len(ss))

	for i := 0; i < len(ss); i++ {
		channels[i] = make(chan FreqMap)
		go channelCount(ss[i], channels[i])
	}
	for _, m := range channels {
		for k, v := range <-m {
			result[k] += v
		}
	}

	return result
}
