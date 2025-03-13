package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

var data = []Entry{}
var MIN = 0
var MAX = 26

func list() {
	for _, v := range data {
		fmt.Println(v)
	}
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func getString(l int64) string {
	startChar := "A"
	temp := ""
	var i int64 = 1
	for {
		myRand := random(MIN, MAX)
		newChar := string(startChar[0] + byte(myRand))
		temp = temp + newChar
		if i == l {
			break
		}
		i++
	}
	return temp
}

func populate(n int, s []Entry) {
	for i := 0; i < n; i++ {
		name := getString(4)
		surname := getString(5)
		n := strconv.Itoa(random(100, 199))
		data = append(data, Entry{name, surname, n, getString(7)})
	}
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Usage: insert|delete|search|list <arguments>")
		return
	}

	if err := SetCSV(); err != nil {
		fmt.Println(err)
	}

	fileInfo, err := os.Stat(CSVFILE)

	if err != nil {
		fmt.Println(err)
		return
	}

	mode := fileInfo.Mode()
	if !mode.IsRegular() {
		fmt.Println(CSVFILE, "not regular file!")
		return
	}

	err = ReadCSVFile(CSVFILE)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = CreateIndex()
	if err != nil {
		fmt.Println("Cannot create index")
	}

	// Differentiate between the commands
	switch arguments[1] {
	case "search":
		if len(arguments) != 3 {
			fmt.Println("Usage: search Number")
			return
		}
		t := strings.ReplaceAll(arguments[2], "-", "")
		if !MatchTel(t) {
			fmt.Println("Not a valid teoephone number:", t)
			return
		}
		temp := Search(t)
		if temp == nil {
			fmt.Println("Number not found:", t)
			return
		}
		fmt.Println(t)
	case "list":
		list()
	default:
		fmt.Println("Not a valid option")
	}
}
