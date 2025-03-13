package main

type Entry struct {
	Name       string
	Surname    string
	Tel        string
	LastAccess string
}

var CSVFILE = "/tmp/csv.data"

var DATA = []Entry{}

var INDEX map[string]int
