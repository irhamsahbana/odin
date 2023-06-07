// Package users is implements component logic.
package users

import (
	"context"
	"reflect"

	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/adapters"
	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/entity"
	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/ports/rest"
	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/usecase"
)

func init() {
	usecase.Register(usecase.Registration{
		Name: "users",
		Inf:  reflect.TypeOf((*T)(nil)).Elem(),
		New: func() any {
			return &impl{}
		},
	})
}

// T is the interface implemented by all users Component implementations.
type T interface {
	Get(ctx context.Context, paging entity.RequestGetUsers) ([]entity.User, rest.Pagination, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
}

type impl struct {
	adapter *adapters.Adapter
}

// Init initializes the execution of a process involved in a users Component usecase.
func (i *impl) Init(adapter *adapters.Adapter) error {
	i.adapter = adapter
	return nil
}
