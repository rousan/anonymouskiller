package main

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rousan/anonymouskiller/stulish"
)

const defaultMessage = "Enjoy the wrath 😁"
const attackDuration = 2 * time.Hour

var pageTemplates *template.Template

func init() {
	tmpl, err := template.ParseGlob("./web/*.html")
	if err != nil {
		log.Fatal(err)
	}
	pageTemplates = tmpl
}

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("PORT env must be set!")
	}

	fs := http.FileServer(http.Dir("./web/static"))
	staticHandler := http.StripPrefix("/static/", fs)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", strings.ToUpper(r.Method), r.RequestURI)

		switch path := r.URL.Path; {
		case path == "/":
			indexGET(w, r)
		case path == "/hack":
			hackGET(w, r)
		case strings.HasPrefix(path, "/static/"):
			staticHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	addr := net.JoinHostPort("0.0.0.0", port)

	log.Printf("Running server on: %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func indexGET(w http.ResponseWriter, r *http.Request) {
	err := pageTemplates.ExecuteTemplate(w, "index.html", "")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func hackGET(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		err = pageTemplates.ExecuteTemplate(w, "error.html", "Invalid Stulish userid or message")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	userid := r.Form.Get("userid")
	message := r.Form.Get("message")

	if userid == "" {
		err = pageTemplates.ExecuteTemplate(w, "error.html", "Stulish User ID is not specified")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if message == "" {
		message = defaultMessage
	}

	c, err := stulish.Attack(userid, message)
	if err != nil {
		err = pageTemplates.ExecuteTemplate(w, "error.html", "Internal Server Error")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	go func() {
		time.Sleep(attackDuration)
		c <- 1
	}()

	err = pageTemplates.ExecuteTemplate(w, "hack.html", "")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
