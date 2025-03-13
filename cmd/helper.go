package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

func SetCSV() error {
	_, err := os.Stat(CSVFILE)
	if err != nil {
		fmt.Println("Creating", CSVFILE)
		f, err := os.Create(CSVFILE)
		if err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return nil
}

func ReadCSVFile(filepath string) error {
	_, err := os.Stat(filepath)
	if err != nil {
		return err
	}
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil
	}
	for _, line := range lines {
		temp := Entry{
			Name:       line[0],
			Surname:    line[1],
			Tel:        line[2],
			LastAccess: line[3],
		}

		DATA = append(DATA, temp)
	}
	return nil
}

func CreateIndex() error {
	index := make(map[string]int)
	for i, k := range DATA {
		key := k.Tel
		index[key] = i
	}
	return nil
}

func MatchTel(s string) bool {
	t := []byte(s)
	re := regexp.MustCompile(`\d+$`)
	return re.Match(t)
}

func Search(key string) *Entry {
	i, ok := INDEX[key]
	if !ok {
		return nil
	}
	DATA[i].LastAccess = strconv.FormatInt(time.Now().Unix(), 10)
	_ = saveCSVFile(CSVFILE)
	return &DATA[i]
}

func saveCSVFile(filepath string) error {
	csvfile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer csvfile.Close()
	csvwrite := csv.NewWriter(csvfile)
	for _, row := range DATA {
		temp := []string{row.Name, row.Surname, row.Tel, row.LastAccess}
		_ = csvwrite.Write(temp)
	}
	csvwrite.Flush()
	return nil
}
