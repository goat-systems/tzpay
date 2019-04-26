package options

import (
	"errors"
	"regexp"
	"strconv"
)

var (
	reCycles = regexp.MustCompile(`([0-9]+)`)
)

// Options is a struct to represent configuration options for payman
type Options struct {
	Delegate        string
	Secret          string
	Password        string
	Service         bool
	Cycles          string
	Cycle           int
	Node            string
	Port            string
	Fee             float32
	File            string
	NetworkFee      int
	NetworkGasLimit int
	Dry             bool
	RedditAgent     string
	RedditTitle     string
}

// ParseCyclesInput parses a string of cycles formated like ParseCyclesInput "10-14"
func (o *Options) ParseCyclesInput() ([2]int, error) {
	arrayCycles := reCycles.FindAllStringSubmatch(o.Cycles, -1)
	if arrayCycles == nil || len(arrayCycles) > 2 {
		return [2]int{}, errors.New("unable to parse cycles flag. Example format 8-12")
	}
	var cycleRange [2]int

	if len(arrayCycles) == 1 {
		cycleRange[0], _ = strconv.Atoi(arrayCycles[0][1])
		cycleRange[1], _ = strconv.Atoi(arrayCycles[0][1])
	} else {
		cycleRange[0], _ = strconv.Atoi(arrayCycles[0][1])
		cycleRange[1], _ = strconv.Atoi(arrayCycles[1][1])
	}

	return cycleRange, nil
}
