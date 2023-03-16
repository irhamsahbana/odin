// Package rest is port handler via http/s protocol
package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"

	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/ports/rest"
	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/shared/entity"
	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/shared/tracer"
	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/usecase/pokemon"
)

// Hello handler instance data.
type Hello struct {
	Pokemon pokemon.T
}

// ResponseWorld - hello handler response.
type ResponseWorld struct {
	Message string `json:"message"`
}

// Register is endpoint group for handler.
func (h *Hello) Register(router chi.Router) {
	router.Get("/hello", rest.JSON(rest.HandlerAdapter[ResponseWorld](h.World)))
	router.Get("/pokemon", rest.JSON(rest.HandlerAdapter[ResponseFetchPokemon](h.FetchPokemon)))
}

// World - endpoint func.
func (*Hello) World(_ http.ResponseWriter, r *http.Request) (ResponseWorld, error) {
	ctx := r.Context()
	log.Error().Err(fmt.Errorf("ini sebuah kesalahan")).Msg("handler World")
	span := trace.SpanFromContext(ctx)
	defer span.End()
	l := log.Hook(tracer.TraceContextHook(ctx))

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

// ResponseFetchPokemon - hello handler response.
type ResponseFetchPokemon struct {
	*entity.Resource
}

// FetchPokemon - endpoint func to get all pokemon.
func (h *Hello) FetchPokemon(w http.ResponseWriter, r *http.Request) (ResponseFetchPokemon, error) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	l := log.Hook(tracer.TraceContextHook(ctx))

	result, err := h.Pokemon.GetAll(ctx)
	if err != nil {
		return ResponseFetchPokemon{}, rest.ErrBadRequest(w, r, err)
	}
	l.Info().Msg("FetchPokemon")
	return ResponseFetchPokemon{result}, nil
}
