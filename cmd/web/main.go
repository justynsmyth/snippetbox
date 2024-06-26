package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// Custom Structured logger dependency needs to be injected to handlers
// We define handler functions as methods against application struct
type application struct {
	logger *slog.Logger
}

// -----------------------------------------------------------------------------
// Purpose:
// Parsing the runtime configuration settings for the application;
// Establishing the dependencies for the handlers;
// Running the HTTP server.
// -----------------------------------------------------------------------------
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

	// Initialize a new instance of our application struct, containing the
	// dependencies (for now, just the structured logger).
	app := &application{
		logger: logger,
	}

	// logger.Info("starting server", "addr", *addr)
	// the below method avoids key/value !BADKEY issues
	// It will force attribute pairs together
	logger.Info("starting server", slog.Any("addr", ":4000"))
	// Call the new app.routes() method to get the servemux containing our routes,
	// and pass that to http.ListenAndServe().
	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
