package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/risor-io/risor"
	"github.com/risor-io/risor/errz"
)

type Request struct {
	Code string `json:"code"`
}

type Response struct {
	Result json.RawMessage `json:"result"`
	Time   float64         `json:"time"`
}

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
	var req Request
	var res Response
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Expected a JSON body in the request", http.StatusBadRequest)
		return
	}

	if req.Code == "" {
		http.Error(w, "Please provide a code snippet", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	start := time.Now()

	result, err := risor.Eval(ctx,
		string(req.Code),
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

	res.Result, err = json.Marshal(result)
	if err != nil {
		http.Error(w, "Unable to marshal output", http.StatusInternalServerError)
		return
	}

	res.Time = time.Since(start).Seconds()

	response, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Unable to marshal output", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
