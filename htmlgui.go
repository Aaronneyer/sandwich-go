package main

import (
	"html/template"
	"net/http"
	"net"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/query_result.html"))

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	f := Search(r.FormValue("search"), r.FormValue("regex")=="true")
	templates.ExecuteTemplate(w, "query_result.html", f)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	DownloadQueue <- &IPFilePair{IP: net.ParseIP(r.FormValue("ip")), FileName: r.FormValue("file")}
}

func killHandler(w http.ResponseWriter, r *http.Request) {
	Shutdown()
}

func InitializeFancyStuff() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/kill", killHandler)
	http.Handle("/static/", http.FileServer(http.Dir("./")))
	http.ListenAndServe(":8000", nil)
}

