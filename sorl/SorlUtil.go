package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

const (
	cnt = 100
)

// Property for maps
type Property map[string]string

func (prop Property) replaceProp(line string) (string, error) {
	return replaceProp(line, prop)
}

func replaceProp(line string, lMap Property) (string, error) {

	sIdx := -1
	lIdx := -1
	lKey := ""
	lMapKey := ""
	lMapVal := ""
	ok := false
	//bakLine := line

	for {
		sIdx = strings.LastIndex(line, "{")
		if sIdx < 0 {
			break
		}
		tLine := line[sIdx+1:]
		lIdx = strings.Index(tLine, "}")
		lKey = line[sIdx+1 : sIdx+lIdx+1]
		lMapKey = strings.TrimSpace(lKey)

		if lMapVal, ok = lMap[lMapKey]; !ok {
			return "", errors.New("Key:'" + lMapKey + "' not found")
		}
		//fmt.Printf("\n%s", lKey)
		line = strings.ReplaceAll(line, "{"+lKey+"}", lMapVal)
		//fmt.Printf("\n%s", line)

	}

	return line, nil
}

// PrintList for print
func PrintList(title string, list []string) {

	fmt.Println()
	fmt.Println("\n" + strings.Repeat("=", cnt))
	fmt.Printf("                    %s\n", title)
	fmt.Println(strings.Repeat("=", cnt))

	for i, str := range list {
		fmt.Printf("%d. %v\n", i+1, str)
	}

	fmt.Println(strings.Repeat("=", cnt))
	fmt.Println()

}

func printMap(title string, myMap SorlMap) {
	PrintProperties(title, Property(myMap))
}

// PrintProperties to print
func PrintProperties(title string, prop Property) {

	keys := make([]string, 0, len(prop))

	for ky := range prop {
		keys = append(keys, ky)
	}

	sort.Strings(keys)
	fmt.Println()
	fmt.Println("\n" + strings.Repeat("=", cnt))
	fmt.Printf("                    %s\n", title)
	fmt.Println(strings.Repeat("=", cnt))

	for i, ky := range keys {
		fmt.Printf("%v. %v=%v\n", i+1, ky, prop[ky])
	}

	fmt.Println(strings.Repeat("=", cnt))
	fmt.Println()
}

// LoadPropertyFile to load
func LoadPropertyFile(fileName, sep string) (Property, error) {

	prp := Property{}

	data, err := ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	for _, str := range data {

		str = strings.TrimLeft(str, " ")
		if strings.HasPrefix(str, "#") {
			continue
		}

		if !strings.Contains(str, "=") {
			continue
		}

		arStr := strings.Split(str, sep)

		if len(arStr) > 0 {
			prp[arStr[0]] = strings.Join(arStr[1:], sep)
		}

	}
	return prp, nil
}

// ReadFile to read a normal file
func ReadFile(fileName string) ([]string, error) {

	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	//return strings.Split(string(data), "\n"), nil

	// Remove EOF char from []string
	lines := strings.Split(string(data), "\n")
	lLen := len(lines)
	return lines[:lLen-1], nil
}

func chkDir(pathName string) bool {

	_, err := os.Stat(pathName)

	if err != nil {
		return false
	}

	return true
}

func chkFileOrDir(pathName string) (bool, error) {

	fi, err := os.Stat(pathName)

	if err != nil {
		return false, err
	}

	switch fdMode := fi.Mode(); {

	case fdMode.IsDir():
		return true, nil
	case fdMode.IsRegular():
		return false, nil

	}

	return false, errors.New("Type Not Available")
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("\nerror: %s\n", err)
		os.Exit(-1)
	}
}
