package api

import (
	"github.com/gorilla/mux"
	"net/http/pprof"
)

func addPprofHandler(r *mux.Router) {
	// http.Handle("/debug/pprof/", http.HandlerFunc(Index))
	// 	http.Handle("/debug/pprof/cmdline", http.HandlerFunc(Cmdline))
	// 	http.Handle("/debug/pprof/profile", http.HandlerFunc(Profile))
	// 	http.Handle("/debug/pprof/symbol", http.HandlerFunc(Symbol))
	// 	http.Handle("/debug/pprof/trace", http.HandlerFunc(Trace))
	r.Path("/pprof/cmdline").Methods("GET").HandlerFunc(pprof.Cmdline)
	r.Path("/pprof/profile").Methods("GET").HandlerFunc(pprof.Profile)
	r.Path("/pprof/symbol").Methods("GET").HandlerFunc(pprof.Symbol)
	r.Path("/pprof/trace").Methods("GET").HandlerFunc(pprof.Trace)
	r.PathPrefix("/debug/pprof").Methods("GET").HandlerFunc(pprof.Index)
}
