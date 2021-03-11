package handlers

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net/http"

	"github.com/cedrickchee/gowebservices/foundation/web"
)

type check struct {
	log *log.Logger
}

func (c check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return errors.New("untrusted error")
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}
	c.log.Println(r, status)
	return web.Respond(ctx, w, status, http.StatusOK)
}
