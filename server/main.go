package main

import (
	"net/http"

	"github.com/po3rin/llb2dot"
)

func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}

// Handler handle request.
func Handler(w http.ResponseWriter, r *http.Request) {
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

	err = llb2dot.WriteDOT(w, g)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(200)
}
