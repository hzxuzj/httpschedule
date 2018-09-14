package api

import (
	_ "fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func addViewHandler(r *mux.Router) {
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts/", http.FileServer(http.Dir("./fonts/"))))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./images/"))))

	r.Path("/").Methods("GET").HandlerFunc(loginHandler)
	r.Path("/viewApp").Methods("GET").HandlerFunc(loginHandler)
	r.Path("/viewImage").Methods("GET").HandlerFunc(loginHandler)
	r.Path("/dashboard").Methods("GET").HandlerFunc(loginHandler)
	r.Path("/viewService").Methods("GET").HandlerFunc(loginHandler)
	r.Path("/viewResource").Methods("GET").HandlerFunc(loginHandler)
	r.Path("/viewVolumes").Methods("GET").HandlerFunc(loginHandler)
	r.Path("/viewNode").Methods("GET").HandlerFunc(loginHandler)
	r.Path("/viewUser").Methods("GET").HandlerFunc(loginHandler)

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/html/index.html") //main.go
	if err != nil {

		log.Println(err)
	}

	t.Execute(w, nil)
}
