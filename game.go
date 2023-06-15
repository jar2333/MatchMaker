package main

const (
	TIE = ""
)

type game interface {
	p1() string
	p2() string

	finished() chan string
}
