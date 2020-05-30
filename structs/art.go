package art

import (
	"fmt"
	"strings"
)

type art struct {
	i   int
	arr [][][]string
}

func (a *art) init() {
	a.arr = append(a.arr, [][]string{})
}

func (a *art) update() {
	a.i++
	a.init()
}

func (a *art) apply(str string, b banner) {
	for i := 0; i < len(str); i++ {
		if str[i] == '\\' && i+1 < len(str) {
			if str[i+1] == 'n' {
				a.update()
				i++
				if i == len(str)-1 {
					break
				}
			} else if str[i+1] == '\'' {
				a.arr[a.i] = append(a.arr[a.i], b.toBig('\''))
				i++
			} else if str[i+1] == '"' {
				a.arr[a.i] = append(a.arr[a.i], b.toBig('"'))
			} else if str[i+1] == '!' {
				a.arr[a.i] = append(a.arr[a.i], b.toBig('!'))
				i++
			} else if str[i+1] == '\\' {
				a.arr[a.i] = append(a.arr[a.i], b.toBig('\\'))
				i++
			} else if str[i+1] == 't' {
				for j := 0; j < 4; j++ {
					a.arr[a.i] = append(a.arr[a.i], b.toBig(' '))
				}
				i++
			}
			continue
		}
		if i+1 < len(str) && str[i] == '\r' && str[i+1] == '\n' {
			a.update()
			i++
			if i == len(str)-1 {
				break
			}
			continue
		}
		big := b.toBig(str[i])
		a.arr[a.i] = append(a.arr[a.i], big)
	}
	if len(a.arr[a.i]) == 0 {
		a.arr[a.i] = append(a.arr[a.i], []string{"", "", "", "", "", "", "", ""})
	}
}

func (a art) simplePrint(index int) {
	for i := 0; i < 8; i++ {
		for j := range a.arr[index] {
			fmt.Print(a.arr[index][j][i])
		}
		fmt.Println()
	}
}

func (a art) Print(b banner) {
	for i := 0; i <= a.i; i++ {
		a.simplePrint(i)
	}
}

func (a art) toStr() string {
	str := ""
	for i := 0; i <= a.i; i++ {
		for j := 0; j < 8; j++ {
			for k := range a.arr[i] {
				str += a.arr[i][k][j]
			}
			str += "\n"
		}
	}
	return str
}

func toArr(str string) [][]string {
	index, newlines := 0, 0
	arr := [][]string{{"", "", "", "", "", "", "", ""}}
	for len(str) > 0 {
		arr[index][newlines%8] += str[:strings.Index(str, "\n")]
		str = str[strings.Index(str, "\n")+1:]
		if newlines%8 == 7 && len(str) > 1 {
			index++
			arr = append(arr, []string{"", "", "", "", "", "", "", ""})
		}
		newlines++
	}
	return arr
}

func AsciiArt(str, font string) string {
	if len(str) == 0 {
		return ""
	}
	a := art{}
	a.init()
	b := banner{}
	b.init("banners/standard.txt")
	if font == "shadow" {
		b.init("banners/shadow.txt")
	} else if font == "thinkertoy" {
		b.init("banners/thinkertoy.txt")
	}

	a.apply(str, b)

	return a.toStr()
}
