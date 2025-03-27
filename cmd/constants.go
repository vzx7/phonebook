package main

type Entry struct {
	Name       string
	Surname    string
	Tel        string
	LastAccess string
}

type Phonebook []Entry

var CSVFILE = "/tmp/csv.data"

var DATA = Phonebook{}

var INDEX map[string]int

var MIN = 0
var MAX = 26
