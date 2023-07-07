package main

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/risor-io/risor"
	"github.com/risor-io/risor/errz"
)

const MaxCodeSize = 100 * 1024

func main() {
	var port string
	flag.StringVar(&port, "port", "8000", "Define port for the server to listen on")
	flag.Parse()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/execute", func(w http.ResponseWriter, r *http.Request) {
		executeHandler(w, r)
	})

	log.Println("Server started on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func executeHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	code, err := io.ReadAll(io.LimitReader(r.Body, MaxCodeSize))
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	result, err := risor.Eval(ctx,
		string(code),
		risor.WithDefaultBuiltins(),
		risor.WithDefaultModules())
	if err != nil {
		if friendlyErr, ok := err.(errz.FriendlyError); ok {
			http.Error(w, friendlyErr.FriendlyErrorMessage(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Unable to marshal result", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
