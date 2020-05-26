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

func indexHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "templates/404.html")
		return
	}

	value := r.FormValue("text")
	font := r.FormValue("font")
	color := r.FormValue("color")

	cmd := exec.Command("./test", value, font)
	cmd.Stdin = os.Stdin
	out, _ := cmd.Output()
	output := string(out)

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
