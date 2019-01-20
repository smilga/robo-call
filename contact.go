package main

const (
	notResponding contactStatus = iota
	responed
)

type contactStatus int

type contact struct {
	number       string
	status       contactStatus
	callDuration int
	inProcess    bool
}
