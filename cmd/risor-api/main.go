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
	"github.com/risor-io/risor/modules/all"
	"github.com/risor-io/risor/parser"
)

type Request struct {
	Content string `json:"content"`
}

type Response struct {
	Result string `json:"result"`
	Time   float64 `json:"time"`
}

func main() {
	var port string
	flag.StringVar(&port, "port", "3000", "Define port for the server to listen on")
	flag.Parse()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/execute", func(w http.ResponseWriter, r *http.Request) {
		executeHandler(w, r)
	})

	log.Println("Server started on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":" + port, r))
}

func executeHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	var res Response
	json.NewDecoder(r.Body).Decode(&req)

	input := req.Content
	if input == "" {
		http.Error(w, "Please provide a code snippet", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	start := time.Now()

	result, err := risor.Eval(ctx, string(input),
		risor.WithBuiltins(all.Builtins()),
		risor.WithDefaultBuiltins(),
		risor.WithDefaultModules())
	if err != nil {
		parserErr, ok := err.(parser.ParserError)
		if ok {
			http.Error(w, parserErr.FriendlyMessage(), http.StatusInternalServerError)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res.Result = result.Inspect()

	res.Time = time.Since(start).Seconds()

	response, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
