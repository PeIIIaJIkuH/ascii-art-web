package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var index = template.Must(template.ParseFiles("templates/index.html"))
var output = template.Must(template.ParseFiles("templates/output.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	index.Execute(w, nil)
}

func outputHandler(w http.ResponseWriter, r *http.Request) {
	value := r.FormValue("text")
	font := r.FormValue("font")
	color := r.FormValue("color")

	cmd := exec.Command("./test", value, font)
	cmd.Stdin = os.Stdin
	out, _ := cmd.Output()
	outp := string(out)

	file, _ := ioutil.ReadFile("templates/output.html")
	str := string(file)

	find := "class=\"output\">"
	size := len(find)
	colorstr := "style=color:" + color + "; "
	str = str[:strings.Index(str, find)] + colorstr + str[strings.Index(str, find):]
	i := strings.Index(str, find) + size
	str = str[:i] + outp + str[i:]

	fmt.Fprint(w, str)

	// output.Execute(w, nil)
}

func main() {
	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "3000"
	// }

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("templates/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/output", outputHandler)

	ip := os.Getenv("APP_IP")
	port := os.Getenv("APP_PORT")

	fmt.Println(ip + ":" + port)

	log.Fatal(http.ListenAndServe(ip+":"+port, mux))
}
