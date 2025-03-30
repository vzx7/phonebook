package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/vzx7/phonebook/pkg/db"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Usage: insert|delete|search|list <arguments>")
		return
	}

	db.Connect([]string{"", "localhost", "5437", "xz", "pass", "master"})

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
			fmt.Println("Not a valid telephone number:", t)
			return
		}
		temp := Search(t)
		if temp == nil {
			fmt.Println("Number not found:", t)
			return
		}
		fmt.Println(*temp)
	case "delete":
		if len(arguments) != 3 {
			fmt.Println("Usage: delete Number")
			return
		}
		t := strings.ReplaceAll(arguments[2], "-", "")
		if !MatchTel(t) {
			fmt.Println("Not a valid telephone number:", t)
			return
		}
		err := DeleteEntry(t)
		if err != nil {
			fmt.Println(err)
		}
	case "insert":
		if len(arguments) != 5 {
			fmt.Println("Usage: insert Name Surname Telephone")
			return
		}
		t := strings.ReplaceAll(arguments[4], "-", "")
		if !MatchTel(t) {
			fmt.Println("Not a valid telephone number:", t)
			return
		}
		temp := InitS(arguments[2], arguments[3], t)
		if temp != nil {
			err = Insert(temp)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	case "list":
		List()
	default:
		fmt.Println("Not a valid option")
	}
}
