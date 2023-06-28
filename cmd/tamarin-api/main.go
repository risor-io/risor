package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cloudcmds/tamarin/v2"
	"github.com/cloudcmds/tamarin/v2/parser"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Request struct {
    ID      int    `json: "id"`
    Content string `json: "content"`
}

func main() {
	r := chi.NewRouter()
    r.Use(middleware.Logger)
	r.Post("/execute", func(w http.ResponseWriter, r *http.Request) {
        executeHandler(w, r)
    })

	log.Println("Server started on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}

func executeHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req Request 
    json.NewDecoder(r.Body).Decode(&req)

	fmt.Println("qwe",req)

	// if err != nil {
	// 	http.Error(w, "Failed to parse request body", http.StatusBadRequest)
	// 	return
	// }

	input := req.Content

	ctx := context.Background()
	start := time.Now()

	result, err := tamarin.Eval(ctx, string(input))
	if err != nil {
		parserErr, ok := err.(parser.ParserError)
		if ok {
			http.Error(w, parserErr.FriendlyMessage(), http.StatusInternalServerError)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := result.Inspect()
	
	response += fmt.Sprintf("\n%.03f", time.Since(start).Seconds())

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, response)
}




// package main

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"time"

// 	"github.com/cloudcmds/tamarin/v2"
// 	"github.com/cloudcmds/tamarin/v2/object"
// 	"github.com/cloudcmds/tamarin/v2/parser"
// )

// func main() {
// 	ctx := context.Background()
// 	start := time.Now()

// 	// data to get
// 	showTiming := true
// 	// input := []byte(`print("sentence")`)

// 	result, err := tamarin.Eval(ctx,
// 		`print("sentence")`,
// 	)
// 	if err != nil {
// 		parserErr, ok := err.(parser.ParserError)
// 		if ok {
// 			fmt.Fprintf(os.Stderr, "%s\n", parserErr.FriendlyMessage())
// 		} else {
// 			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
// 		}
// 		os.Exit(1)
// 	}
// 	if result != object.Nil {
// 		fmt.Println(result.Inspect())
// 	}
// 	if showTiming {
// 		fmt.Printf("%.03f\n", time.Since(start).Seconds())
// 	}

// }