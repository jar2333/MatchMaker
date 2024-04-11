package utils

import (
	"bufio"
	"os"
	"strings"
	"time"
)

func LoadDate() time.Time {
	dat, err := os.ReadFile("./date.txt")
	if err != nil {
		panic(err)
	}

	// Calling Parse() method with its parameters
	tm, e := time.Parse(time.RFC822, string(dat))
	if e != nil {
		panic(e)
	}

	return tm
}

func LoadKeys() map[string]bool {
	keys := make(map[string]bool)

	dat, err := os.ReadFile("./keys.txt")
	if err != nil {
		panic(err)
	}

	lines := splitLines(string(dat))

	for _, l := range lines {
		keys[l] = true
	}

	return keys
}

func splitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}
