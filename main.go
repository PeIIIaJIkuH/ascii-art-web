package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
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

	cmd := exec.Command("./test", value, font)
	cmd.Stdin = os.Stdin
	out, _ := cmd.Output()
	output := string(out)

	if len(output) == 0 {
		temp := exec.Command("test", value, font)
		temp.Stdin = os.Stdin
		out1, _ := temp.Output()
		output1 := string(out1)

		fmt.Fprintln(w, output1)
	}

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
