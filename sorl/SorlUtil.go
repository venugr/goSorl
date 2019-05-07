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

// Property
type Property map[string]string

// PrintList
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

func ReadFile(fileName string) ([]string, error) {

	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	return strings.Split(string(data), "\n"), nil

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
