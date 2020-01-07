package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

type Section struct {
	str   string
	count int
}

// String() : simple format output for Section struct
func (t Section) String() string {
	return fmt.Sprintf("%d:%s", t.count, t.str)
}

// TopN : sorts the vals slice in descending order ,then takes
// the top 5 entries and returns them
func TopN(vals []Section) []Section {
	// sorts vals in descending order of count
	sort.Slice(vals, func(i, j int) bool {
		return vals[i].count > vals[j].count
	})

	//takes top 5 items from vals
	if cap(vals) < 5 {
		return vals
	} else {
		return vals[:5]
	}
}

func findRepeats(buf []byte, repeatLen int, minCount int) []Section {
	m := map[string]Section{}
	size := len(buf)

	// make a map of all substrings of repeatLen to occurrances
	for left := 0; left < size-repeatLen; left++ {

		//get string of length repeatLen at current position
		current := string(buf[left : left+repeatLen])

		// if already in map, increment
		// else, add to map
		if value, exist := m[current]; exist {
			value.count++
			m[current] = value
		} else {
			m[current] = Section{current, 1}
		}
	}

	// make a slice of values with min count
	vals := []Section{}

	// iterate through map and add items with at least
	// minCount occurrences to vals
	for k, v := range m {
		if m[k].count >= minCount {
			vals = append(vals, v)
		}
	}

	// take topN and return
	vals = TopN(vals)
	return vals
}

// addValueIfNotSubstr : Add items in vals slice to the results
// only if they are not substrings of any items in results.
// returns results with appended values.
// Reason: If GGGTTTAAA is in results, then we don't need GGGTTT
// or TTTAAA, which are substrings.
func addValueIfNotSubstr(results []Section, vals []Section) []Section {
	itemsToAdd := []Section{}
	for _, val := range vals {
		// is it a substring of previous?
		skip := false
		for _, prev := range results {
			if strings.Contains(prev.str, val.str) {
				delta := math.Abs(float64(prev.count - val.count))
				if delta <= 2 {
					//log.Printf("%s(%d) contains %s(%d)\n", prev.str, prev.count, val.str, val.count)
					skip = true
				}
				break
			}
		}
		if !skip {
			itemsToAdd = append(itemsToAdd, val)
		}
	}

	return append(results, itemsToAdd...)
}

// minCountForResults is the minimum times a section occurs that
// we will include in results
var minCountForResults = 12
var minRepeatLen = 25
var maxRepeatLen = 35
var nucleotideBufferSize = 2200

func processBuffer(buf []byte, lineStart int, lineEnd int) {
	//log.Println("processBuffer line:", lineNum)
	results := []Section{}
	for i := maxRepeatLen; i > minRepeatLen; i-- {
		vals := findRepeats(buf, i, minCountForResults)
		results = addValueIfNotSubstr(results, vals)
	}
	if len(results) > 0 {
		fmt.Println("lines:", lineStart, "-", lineEnd, results)
	}
}

func handleCommandLineFlags() {
	pMinMatchFlag := flag.Int("minmatches", minCountForResults, "minimum number of occurrances to include in results")
	pMinRepeatLenFlag := flag.Int("minrepeatlen", minRepeatLen, "minimum length of repeating section")
	pMaxRepeatLenFlag := flag.Int("maxrepeatlen", maxRepeatLen, "maximum length of repeating section")
	pBufsizeFlag := flag.Int("bufsize", nucleotideBufferSize, "number of nucleotides in analysis buffer")

	flag.Parse()

	minCountForResults = *pMinMatchFlag
	minRepeatLen = *pMinRepeatLenFlag
	maxRepeatLen = *pMaxRepeatLenFlag
	nucleotideBufferSize = *pBufsizeFlag
}

func main() {
	numLines := 0
	lineStart := 0

	handleCommandLineFlags()

	buf := []byte{}
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Bytes()

		// skip FASTA header line and lines with NNN...
		if line[0] == byte('>') || strings.Contains(string(line), "NNNNNN") {
			numLines++
			continue
		}

		buf = append(buf, line...)

		if len(buf) >= nucleotideBufferSize {
			processBuffer(buf, lineStart, numLines)
			// clear out buffer
			buf = []byte{}
			lineStart = numLines
		}
		numLines++
	}
}
