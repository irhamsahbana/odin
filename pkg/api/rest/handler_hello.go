// Package rest is port handler via http/s protocol
package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kubuskotak/valkyrie"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"

	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/ports/rest"
)

// Hello handler instance data.
type Hello struct {
}

// ResponseWorld - hello handler response.
type ResponseWorld struct {
	Message string `json:"message"`
}

// Register is endpoint group for handler.
func (h *Hello) Register(router chi.Router) {
	router.Get("/hello", rest.JSON(rest.HandlerAdapter[ResponseWorld](h.World)))
}

// World - endpoint func.
func (*Hello) World(_ http.ResponseWriter, r *http.Request) (ResponseWorld, error) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()
	l := log.Hook(valkyrie.TraceContextHook(ctx))

	l.Info().Str("Hello", "World").Msg("this")
	rest.Paging(r, rest.Pagination{
		Page:  1,
		Limit: 10,
		Size:  100,
		Total: 500,
	})
	rest.StatusCreated(r)
	return ResponseWorld{
		Message: "Hello everybody",
	}, nil
}
