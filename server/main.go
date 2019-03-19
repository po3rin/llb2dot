package main

import (
	"net/http"

	"github.com/po3rin/llb2dot"
)

func main() {
	fs := http.FileServer(http.Dir("view"))
	http.Handle("/", fs)
	http.HandleFunc("/api/dot", Handler)
	http.ListenAndServe(":8080", nil)
}

// Handler handle request.
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	ops, err := llb2dot.LoadDockerfile(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	g, err := llb2dot.LLB2Graph(ops)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(200)
	err = llb2dot.WriteDOT(w, g)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
