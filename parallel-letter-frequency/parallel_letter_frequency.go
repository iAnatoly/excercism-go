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

// naiveConcurrentFrequency is a aaive implementation w/o concurrency.
// Not used anywhere, created just for reference.
// Curiosly enough, this naive implementation passes the benchmarks.
func naiveConcurrentFrequency(ss []string) FreqMap {
	m := FreqMap{}
	for _, s := range ss {
		for _, r := range s {
			m[r]++
		}
	}
	return m
}

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

// ConcurrentFrequency counts the letter frequency on an array of strings
// Each elemet in the arry of strings is processed concurrently
// ConcurrentFrequency returns a map of letters (runes) to their respective frequency
func ConcurrentFrequency(ss []string) FreqMap {
	smap := safeMap{freqmap: make(FreqMap)}

	for _, s := range ss {
		smap.wait.Add(1)
		go safeCount(s, &smap)
	}
	smap.wait.Wait() // wait for all SafeCount threads to finish
	return smap.freqmap
}
