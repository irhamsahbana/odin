// Package rest is port handler via http/s protocol
package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/entity"
	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/ports/rest"
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
	router.Get("/hello", rest.HandlerAdapter[ResponseWorld](h.World).JSON)
	router.Get("/pokemon", rest.HandlerAdapter[ResponseFetchPokemon](h.FetchPokemon).JSON)
}

// World - endpoint func.
func (*Hello) World(_ http.ResponseWriter, r *http.Request) (ResponseWorld, error) {
	_, span, l := tracer.StartSpanLogTrace(r.Context(), "{{$Handler}}")
	defer span.End()

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
	ctx, span, l := tracer.StartSpanLogTrace(r.Context(), "{{$Handler}}")
	defer span.End()

	result, err := h.Pokemon.GetAll(ctx)

	if err != nil {
		return ResponseFetchPokemon{}, rest.ErrBadRequest(w, r, err)
	}
	l.Info().Msg("FetchPokemon")
	return ResponseFetchPokemon{result}, nil
}
