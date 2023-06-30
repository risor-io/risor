// This package implements a Risor language server.
// This is very much a work in progress. The VSCode extension is
// usable as-is but only for syntax highlighting. The language server
// does not currently install when the extension is installed.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jdbaldry/go-language-server-protocol/jsonrpc2"
	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	name := "risor-language-server"
	version := "dev"

	logFile, err := os.OpenFile(
		fmt.Sprintf("%s.log", name),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0666,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logFile.Close()

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(logFile)

	ctx := context.Background()
	stream := jsonrpc2.NewHeaderStream(NewDefaultStdio())
	conn := jsonrpc2.NewConn(stream)
	client := protocol.ClientDispatcher(conn)

	s := Server{
		name:    name,
		version: version,
		client:  client,
		cache:   newCache(),
	}

	conn.Go(ctx, protocol.Handlers(
		protocol.ServerHandler(&s, jsonrpc2.MethodNotFound),
	))
	<-conn.Done()

	if err := conn.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
