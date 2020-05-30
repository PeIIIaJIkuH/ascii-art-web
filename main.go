package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	art "github.com/peiiiajikuh/ascii-art-web/structs"
)

type output struct {
	HasOutput bool
	Text      string
	Color     string
}

func checkValue(str string) bool {
	alpha := art.Alphabet()
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
	exportFormat := r.FormValue("output_format")

	if !checkValue(value) {
		w.WriteHeader(http.StatusBadRequest)
		http.ServeFile(w, r, "templates/400.html")
		return
	}

	data := output{true, art.AsciiArt(value, font), color}

	if exportFormat == "txt" {
		w.Header().Set("Content-Disposition", "attachement; filename=output.txt")
		w.Header().Set("Content-Length", strconv.Itoa(len(data.Text)))
		http.ServeContent(w, r, "output.txt", time.Now(), bytes.NewReader([]byte(data.Text)))
	} else if exportFormat == "pdf" {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "templates/500.html")
		return
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	t.ExecuteTemplate(w, "index", data)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("templates/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/", indexHandler)

	fmt.Println(":" + port)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
