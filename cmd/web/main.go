package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Initialize a new structured logger, which writes to stdout stream w default settings (nil)
	// we can replace slog.NewTextHandler with slog.NewJSONHander() if you want to output in JSON format
	/*
		Any log entries with severity less than Info (Debug) is silently discarded by default
		slog.HandlerOptions{} is used to override this, for example (where nil is below)
		&slog.HandlerOptions{ Level: slog.LevelDebug, }

		We can innclude filename and line number of calling source code
		&slog.HandlerOptions{ AddSource: true, }

		If we want to save log output, in cli we can just redirect stdout to a file
		Ie: go run ./cmd/web >> /tmp/web.log (this appends, not truncates)

		slog.New() is concurrency-safe so we can share it across multiple goroutines in HTTP handlers
	*/
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mux := http.NewServeMux()

	// serves files out of a directory
	fileServer := http.FileServer(http.Dir("./ui/static"))
	// This will find any URL paths that start with "/static/" for matching paths
	// We need to strip /static part as well
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	// logger.Info("starting server", "addr", *addr)
	// the below method avoids key/value !BADKEY issues
	// It will force attribute pairs together
	logger.Info("starting server", slog.Any("addr", ":4000"))

	err := http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
