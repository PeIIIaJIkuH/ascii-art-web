package internal

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var (
	templates      *template.Template
	AsciiArtResult string
)

func init() {
	templates = template.Must(template.ParseGlob("./templates/*.html"))
}

func IndexPage(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.Error(w, "Go back to the main page", 404)
		return
	}
	switch req.Method {
	case "GET":
		indexPageGetHandler(w, req)
	case "POST":
		indexPagePostHandler(w, req)
	default:
		http.Error(w, "Go back to the main page", 405)
		return
	}

}

func indexPageGetHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")

	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "Go back to the main page", 500)
		return
	}
}

func indexPagePostHandler(w http.ResponseWriter, req *http.Request) {

	if errForm := req.ParseForm(); errForm != nil {
		http.Error(w, "Go back to the main page", 400)
		return
	}
	exportFormat := req.FormValue("exportFormat")
	inputText := req.Form.Get("inputText")
	fontFile := req.Form.Get("bannerSelectControl")
	if !isPrintableChar(inputText) {
		http.Error(w, "Go back to the main page", 400)
		return
	}
	if !checkFs(fontFile) {
		http.Error(w, "Go back to the main page", 500)
		return
	}
	AsciiArtResult, errExec := getOutputText(inputText, fontFile)
	if errExec != nil {
		http.Error(w, "Go back to the main page", 500)
		return
	}

	data := Data{
		HasOutput:  true,
		OutputText: AsciiArtResult,
	}
	if exportFormat == "txt" {
		w.Header().Set("Content-Disposition", "attachment; filename=output.txt")
		w.Header().Set("Content-Length", strconv.Itoa(len(data.OutputText)))
		http.ServeContent(w, req, "output.txt", time.Now(), bytes.NewReader([]byte(data.OutputText)))

	} else if exportFormat == "pdf" {
		http.Error(w, "Go back to the main page", 500)
		return
	}
	errTemplate := templates.ExecuteTemplate(w, "index.html", data)
	if errTemplate != nil {
		http.Error(w, "Go back to the main page", 500)
		return
	}
}
