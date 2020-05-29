package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	os.Args[1] = strings.ReplaceAll(os.Args[1], string([]byte{13, 10}), string([]byte{92, 110}))
	printWords(getCharArray(os.Args[2]), getSplitWord(os.Args[1]))
}

func getSplitWord(word string) []string {
	splitWord := strings.Split(word, "\\n")
	return splitWord
}

func getCharArray(filename string) []string {

	contentByte, err := ioutil.ReadFile("./static/banners/" + filename)
	// handle IOError
	if err != nil {
		fmt.Printf("open %s: no such file or directory", filename)
		os.Exit(0)
	}
	strContent := string(contentByte)
	charArr := strings.Split(strContent, "\n")
	return charArr
}

func printWords(charArr []string, splitWord []string) {
	for _, word := range splitWord {
		printChar(word, charArr)
	}
}

func printChar(word string, charArr []string) {
	if word == "" {
		fmt.Print("\n\n\n\n\n\n\n\n")
		return
	}
	for i := 1; i <= 8; i++ {
		for _, char := range word {
			index := int(rune(char)-32) * 9 // starting index in array
			fmt.Print(charArr[index+i])
		}
		fmt.Println()
	}
}
