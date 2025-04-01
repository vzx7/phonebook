package main

type Entry struct {
	Name       string
	Surname    string
	Tel        string
	LastAccess string
}

type Phonebook []Entry

var INDEX map[string]int
var DATA = Phonebook{}
var CSVFILE = "/tmp/csv.data"

const (
	MIN = 0
	MAX = 26
)
