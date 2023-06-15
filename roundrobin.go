package main

const (
	EMPTY_KEY = ""
)

type pair struct {
	p1 string
	p2 string
}

func getSchedule(registered []string) [][]pair {
	keys := make([]string, 0, len(registered)+1)

	copy(keys, registered)

	if len(keys)%2 == 1 {
		keys = append(keys, EMPTY_KEY)
	}

	schedule := make([][]pair, 0, len(keys)-1)

	for i := 0; i < len(keys)-1; i++ {
		schedule = append(schedule, getRound(keys))
		shift(keys)
	}

	return schedule
}

func getRound(keys []string) []pair {
	pair_amount := len(keys) / 2

	round := make([]pair, pair_amount)

	for i := 0; i < pair_amount; i++ {
		round[i].p1 = keys[i]
		round[i].p2 = keys[mod(-(i+1), len(keys))]
	}

	return round
}

func shift(keys []string) {
	last := keys[len(keys)-1]

	for i := len(keys) - 1; i > 0; i-- {
		keys[i] = keys[i-1]
	}

	keys[0] = last
}

func mod(i int, n int) int {
	return ((i % n) + n) % n
}
