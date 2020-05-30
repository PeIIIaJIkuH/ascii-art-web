package art

import (
	"io/ioutil"
	"log"
	"strings"
)

type banner struct {
	arr [][]string
}

func (b *banner) clear() {
	b.arr = [][]string{}
}

func (b *banner) init(filename string) {
	b.clear()
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		return
	}
	str := string(file)
	for i := 0; i < 96; i++ {
		b.arr = append(b.arr, []string{})
		str = str[strings.Index(str, "\n")+1:]
		b.arr[i] = append(b.arr[i], []string{"", "", "", "", "", "", "", ""}...)
		for j := 0; j < 8; j++ {
			b.arr[i][j] += str[:strings.Index(str, "\n")]
			str = str[strings.Index(str, "\n")+1:]
		}
	}
}

func isEqual(a1, a2 []string) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i := range a1 {
		if a1[i] != a2[i] {
			return false
		}
	}
	return true
}

func (b banner) Index(symbol []string) int {
	for i, j := range b.arr {
		if isEqual(symbol, j) {
			return i
		}
	}
	return -1
}

func (b banner) toBig(symbol byte) []string {
	alpha := Alphabet()
	index := strings.Index(alpha, string(symbol))
	return b.arr[index]
}

func (b banner) Find(big []string) int {
	for i, j := range b.arr {
		if isEqual(big, j) {
			return i
		}
	}
	return -1
}

func Alphabet() string {
	str := ""
	for i := 32; i <= 126; i++ {
		str += string(rune(i))
	}
	return str
}
