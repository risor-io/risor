// main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/risor-io/risor"
	rsql "github.com/risor-io/risor/modules/sql"
)

func main() {

	script, err := os.ReadFile("example.risor")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	wrappedDB := rsql.NewFromDB(context.Background(), db)

	// Multiple evals can reuse the same connection
	for i := 0; i < 3; i++ {
		fmt.Printf("\n# Run %d\n\n", i+1)
		_, err = risor.Eval(context.Background(), string(script), risor.WithGlobal("db", wrappedDB))
		if err != nil {
			log.Fatal(err)
		}
	}
}
