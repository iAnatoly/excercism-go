package letter

import (
	"github.com/tjarratt/babble"
	"os"
	"reflect"
	"strings"
	"testing"
)

// In the separate file frequency.go, you are given a function, Frequency(),
// to sequentially count letter frequencies in a single text.
//
// Perform this exercise on parallelism using Go concurrency features.
// Make concurrent calls to Frequency and combine results to obtain the answer.

const numDict int = 16
const dictSize int = 1000

var randomLetters []string
var randomLettersJoined string

func TestMain(m *testing.M) {
	babbler := babble.NewBabbler()
	babbler.Separator = " "
	babbler.Count = dictSize

	randomLetters = make([]string, dictSize)

	for i := 0; i < numDict; i++ {
		randomLetters[i] = babbler.Babble()
	}

	randomLettersJoined = strings.Join(randomLetters[:], "")

	os.Exit(m.Run())
}

func OriginalFrequency(s string) FreqMap {
	m := FreqMap{}
	for _, r := range s {
		m[r]++
	}
	return m
}

func TestConcurrentFrequency(t *testing.T) {
	seq := OriginalFrequency(randomLettersJoined)
	con := ConcurrentFrequency(randomLetters)
	if !reflect.DeepEqual(con, seq) {
		t.Fatal("ConcurrentFrequency wrong result")
	}
}

func TestSequentialFrequency(t *testing.T) {
	oSeq := OriginalFrequency(randomLettersJoined)
	seq := Frequency(randomLettersJoined)
	if !reflect.DeepEqual(oSeq, seq) {
		t.Fatal("Frequency wrong result")
	}
}

func BenchmarkSequentialFrequency(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Frequency(randomLettersJoined)
	}
}

func BenchmarkConcurrentLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ConcurrentFrequency(randomLetters)
	}
}

func BenchmarkMutexConcurrentLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MutexConcurrentFrequency(randomLetters)
	}
}
func BenchmarkMultipleChannelsConcurrentLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MultipleChannelsConcurrentFrequency(randomLetters)
	}
}
