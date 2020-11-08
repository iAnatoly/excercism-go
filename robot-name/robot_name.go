package robotname

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Robot type
type Robot struct {
	name string
}

const (
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	maxnum   = 1000
)

var (
	namecache []string
	once      sync.Once
)

func populateMap() {
	for _, l1 := range alphabet {
		for _, l2 := range alphabet {
			for i := 0; i < maxnum; i++ {
				namecache = append(namecache, fmt.Sprintf("%c%c%03d", l1, l2, i))
			}
		}
	}
}

func pickRandomName() (string, error) {
	namespaceLen := len(namecache)
	if namespaceLen == 0 {
		return "", errors.New("namespace exhausted")
	}
	index := rand.Intn(namespaceLen)
	name := namecache[index]
	// deleting slice element. Why is this not a single delte oiprator in Golang?
	namecache[index] = namecache[len(namecache)-1]
	namecache = namecache[:len(namecache)-1]
	return name, nil
}

// Name returns a robotr name, or generates one if not set
func (r *Robot) Name() (string, error) {
	// one-time initialization
	once.Do(func() {
		rand.Seed(time.Now().UTC().UnixNano())
		populateMap()
	})

	// pick a new name if needed
	if r.name == "" {
		name, err := pickRandomName()
		if err != nil {
			return "", err
		}
		r.name = name
	}
	return r.name, nil
}

// Reset returns robot to factory defauilts
func (r *Robot) Reset() {
	r.name = ""
}
