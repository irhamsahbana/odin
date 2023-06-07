// Package rest is port handler.
package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/entity"
	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/ports/rest"
	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/shared/tracer"
	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/usecase/users"
)

// MongoRest handler instance data.
type MongoRest struct {
	UsersUsecase users.T
}

// Register is endpoint group for handler.
func (h *MongoRest) Register(router chi.Router) {
	router.Get("/users", rest.HandlerAdapter[[]entity.User](h.GetAll).JSON)
	router.Post("/user", rest.HandlerAdapter[entity.User](h.Create).JSON)
}

// GetAll user.
func (h *MongoRest) GetAll(w http.ResponseWriter, r *http.Request) ([]entity.User, error) {
	var (
		request entity.RequestGetUsers
	)

	ctx, span, l := tracer.StartSpanLogTrace(r.Context(), "{{$Handler}}")
	defer span.End()

	b, err := rest.Bind(r, &request)
	if err != nil {
		return nil, rest.ErrBadRequest(w, r, err)
	}
	if err := b.Validate(); err != nil {
		return nil, rest.ErrBadRequest(w, r, err)
	}

	documents, paging, err := h.UsersUsecase.Get(ctx, request)
	if err != nil {
		return nil, rest.ErrBadRequest(w, r, err)
	}

	rest.Paging(r, paging)

	l.Info().Msg("GetUsers")
	return documents, nil
}

// Create user.
func (h *MongoRest) Create(w http.ResponseWriter, r *http.Request) (entity.User, error) {
	ctx, span, l := tracer.StartSpanLogTrace(r.Context(), "{{$Handler}}")
	defer span.End()

	request := entity.User{}
	b, err := rest.Bind[entity.User](r, &request)
	if err != nil {
		return entity.User{}, rest.ErrBadRequest(w, r, err)
	}
	if err = b.Validate(); err != nil {
		return entity.User{}, rest.ErrBadRequest(w, r, err)
	}

	documents, err := h.UsersUsecase.Create(ctx, request)
	if err != nil {
		return entity.User{}, rest.ErrBadRequest(w, r, err)
	}

	l.Info().Msg("CreateUser")
	return documents, nil
}
