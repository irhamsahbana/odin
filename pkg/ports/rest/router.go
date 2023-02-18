// Package rest is port adapter via http/s protocol
// # This manifest was generated by ymir. DO NOT EDIT.
package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/infrastructure"
	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/shared/tracer"
	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/version"
)

// RegisterFunc is type func to register handler.
type RegisterFunc func(h chi.Router) http.Handler

// Router is the data struct.
type Router struct {
	h *chi.Mux
}

// Register will assign rest handler.
func (r *Router) Register(fn RegisterFunc) http.Handler {
	// set router middleware
	r.h.Use(tracer.Middleware(
		infrastructure.Envs.App.ServiceName))
	r.h.Use(SemanticVersion(
		infrastructure.Envs.App.ServiceName,
		version.GetVersion().VersionNumber(),
	))
	return fn(r.h)
}

// Routes create Router instance.
func Routes() *Router {
	r := &Router{
		h: chi.NewRouter(), // port http
	}
	return r
}