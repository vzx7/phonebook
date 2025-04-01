package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

func SetCSV() error {
	filepath := os.Getenv("PHONEBOOK")
	if filepath != "" {
		CSVFILE = filepath
	}

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

	fileinfo, err := os.Stat(CSVFILE)
	if err != nil {
		return err
	}

	mode := fileinfo.Mode()
	if !mode.IsRegular() {
		return fmt.Errorf("%s nat a regular file", CSVFILE)
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
	INDEX = make(map[string]int)
	for i, k := range DATA {
		key := k.Tel
		INDEX[key] = i
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

func DeleteEntry(key string) error {
	i, ok := INDEX[key]
	if !ok {
		return fmt.Errorf("%s cannot be found", key)
	}
	DATA = append(DATA[:i], DATA[i+1:]...)
	delete(INDEX, key)
	err := saveCSVFile(CSVFILE)
	if err != nil {
		return err
	}
	return nil
}

func InitS(N, S, T string) *Entry {
	if T == "" || S == "" {
		return nil
	}
	LastAccess := strconv.FormatInt(time.Now().Unix(), 10)
	return &Entry{Name: N, Surname: S, Tel: T, LastAccess: LastAccess}
}

func Insert(pS *Entry) error {
	_, ok := INDEX[(*pS).Tel]
	if ok {
		return fmt.Errorf("%s allredy exist", pS.Tel)
	}
	DATA = append(DATA, *pS)
	_ = CreateIndex()
	err := saveCSVFile(CSVFILE)
	if err != nil {
		return err
	}
	return nil
}

func List() {
	sort.Sort(Phonebook(DATA))
	for _, v := range DATA {
		fmt.Println(v)
	}
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

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandomString(lenght uint64) string {
	startChar := "A"
	rs := ""
	for i := uint64(0); i < lenght; i++ {
		myRand := random(MIN, MAX)
		newChar := string(startChar[0] + byte(myRand))
		rs += newChar
	}
	return rs
}

func (p Phonebook) Len() int {
	return len(p)
}

func (p Phonebook) Less(i, j int) bool {
	if p[i].Surname == p[j].Surname {
		return p[i].Name < p[j].Name
	}
	return p[i].Surname < p[j].Surname
}

func (p Phonebook) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
