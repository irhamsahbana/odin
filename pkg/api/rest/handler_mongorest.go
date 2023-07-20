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
	router.Get("/users", rest.HandlerAdapter(h.GetAll).JSON)
	router.Post("/user", rest.HandlerAdapter(h.Create).JSON)
	router.Get("/testing-kafka", rest.HandlerAdapter(h.KafkaTesting).JSON)
}

// GetAll user.
func (h *MongoRest) GetAll(w http.ResponseWriter, r *http.Request) ([]entity.User, error) {
	ctx, span, l := tracer.StartSpanLogTrace(r.Context(), "GetAll")
	defer span.End()
	var (
		request entity.RequestGetUsers
	)

	b, err := rest.Bind(r, &request)
	if err != nil {
		return nil, rest.ErrBadRequest(w, r, err)
	}
	if err := b.Validate(); err != nil {
		return nil, rest.ErrBadRequest(w, r, err)
	}

	documents, err := h.UsersUsecase.Get(ctx, request)
	if err != nil {
		return nil, rest.ErrBadRequest(w, r, err)
	}

	rest.Paging(r, rest.Pagination{
		Page:  documents.Page,
		Limit: documents.Limit,
	})

	l.Info().Msg("GetUsers")
	return documents.Users, nil
}

// Create user.
func (h *MongoRest) Create(w http.ResponseWriter, r *http.Request) (entity.User, error) {
	ctx, span, l := tracer.StartSpanLogTrace(r.Context(), "Create")
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

func (h *MongoRest) KafkaTesting(w http.ResponseWriter, r *http.Request) (e entity.User, err error) {
	ctx, span, l := tracer.StartSpanLogTrace(r.Context(), "Update")
	defer span.End()

	err = h.UsersUsecase.Update(ctx)
	if err != nil {
		return e, rest.ErrBadRequest(w, r, err)
	}

	l.Info().Msg("UpdateUser")
	return e, nil
}
