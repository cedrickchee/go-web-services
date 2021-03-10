// Package web contains a small web framework extension.
package web

import (
	"context"
	"net/http"
	"os"
	"syscall"

	"github.com/dimfeld/httptreemux/v5"
)

// A Handler is a type that handles an http request within our own little mini
// framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
}

// NewApp creates an App value that handle a set of routes for the application.
// It is a factory function.
func NewApp(shutdown chan os.Signal) *App {
	// Construct an App and initialize the ContextMux
	app := App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
	}

	return &app
}

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// Handle ...
func (a *App) Handle(method string, path string, handler Handler) {
	h := func(w http.ResponseWriter, r *http.Request) {

		// BOILERPLATE

		if err := handler(r.Context(), w, r); err != nil {
			a.SignalShutdown()
			return
		}

		// BOILERPLATE
	}

	a.ContextMux.Handle(method, path, h)
}
