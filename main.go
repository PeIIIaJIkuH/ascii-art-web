package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func alphabet() string {
	str := ""
	for i := 32; i <= 126; i++ {
		str += string(rune(i))
	}
	return str
}

func checkValue(str string) bool {
	alpha := alphabet()
	for i := range str {
		if str[i] != '\n' && str[i] != '\r' && !strings.Contains(alpha, string(str[i])) {
			return false
		}
	}
	return true
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "templates/404.html")
		return
	}

	value := r.FormValue("text")
	font := r.FormValue("font")
	if font == "" {
		font = "standard"
	}
	color := r.FormValue("color")
	if color == "" {
		color = "#ffffff"
	}

	fmt.Println(font)

	if _, err := os.Stat(font + ".txt"); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "templates/500.html")
		return
	}

	if !checkValue(value) {
		w.WriteHeader(http.StatusBadRequest)
		http.ServeFile(w, r, "templates/400.html")
		return
	}

	output := asciiArt(value, font)

	file, _ := ioutil.ReadFile("templates/index.html")
	str := string(file)

	find := "class=\"output\">"
	size := len(find)
	colorstr := "style=color:" + color + "; "
	str = str[:strings.Index(str, find)] + colorstr + str[strings.Index(str, find):]
	i := strings.Index(str, find) + size
	if value != "" {
		str = str[:i] + output + str[i:]
	}

	fmt.Fprint(w, str)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("templates/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/", indexHandler)

	fmt.Println(":" + port)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}

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
	alpha := alphabet()
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

func generateRgb(color string) string {
	switch color {
	case "":
		return "255;255;255m"
	case "white":
		return "255;255;255m"
	case "black":
		return "0;0;0m"
	case "red":
		return "255;0;0m"
	case "green":
		return "0;128;0m"
	case "yellow":
		return "255;0255;0m"
	case "blue":
		return "0;0;255m"
	case "magenta":
		return "255;0;255m"
	case "cyan":
		return "0;255;255m"
	case "lime":
		return "0;255;0m"
	case "silver":
		return "192;192;192m"
	case "gray":
		return "128;128;128m"
	case "maroon":
		return "128;0;0m"
	case "olive":
		return "128;128;0m"
	case "purple":
		return "128;0;128m"
	case "teal":
		return "0;128;128m"
	case "mint":
		return "170;255;195m"
	case "lavender":
		return "230;190;255m"
	case "pink":
		return "250;190;190m"
	case "brown":
		return "170;110;40m"
	case "orange":
		return "245;130;48m"
	case "apricot":
		return "255;215;180m"
	case "beige":
		return "255;250;200m"
	case "tomato":
		return "255;99;71m"
	case "gold":
		return "255;215;0m"
	case "salmon":
		return "250;128;114m"
	default:
		arr := strings.Split(color, ".")
		if len(arr) == 3 {
			r, e1 := strconv.Atoi(arr[0])
			g, e2 := strconv.Atoi(arr[1])
			b, e3 := strconv.Atoi(arr[2])
			if e1 != nil || e2 != nil || e3 != nil || r < 0 || g < 0 || b < 0 || r > 255 || g > 255 || b > 255 {
				return ""
			}
			return strconv.Itoa(r) + ";" + strconv.Itoa(g) + ";" + strconv.Itoa(b) + "m"
		}
		return ""
	}
}

func printStr(str string, count int) {
	for i := 0; i < count; i++ {
		fmt.Print(str)
	}
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

func asciiArt(str, font string) string {
	if len(str) == 0 {
		return "String must contain at least 1 character!"
	}
	a := art{}
	a.init()
	b := banner{}
	b.init("standard.txt")
	if font == "shadow" {
		b.init("shadow.txt")
	} else if font == "thinkertoy" {
		b.init("thinkertoy.txt")
	}

	a.apply(str, b)

	return a.toStr()
}
